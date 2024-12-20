// Copyright 2024 anza-labs contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	anzalabsdevv1alpha1 "github.com/anza-labs/image-builder/api/v1alpha1"
	"github.com/anza-labs/image-builder/internal/naming"
	"github.com/anza-labs/image-builder/internal/s3"
	"github.com/anza-labs/image-builder/pkg/builder"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	// name of image custom finalizer.
	imageFinalizer = "image-builder.anza-labs.dev/finalizer"
)

// ImageReconciler reconciles a Image object.
type ImageReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Builder *builder.Builder
}

// +kubebuilder:rbac:groups=image-builder.anza-labs.dev,resources=images,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=image-builder.anza-labs.dev,resources=images/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=image-builder.anza-labs.dev,resources=images/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ImageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx, "image", klog.KRef(req.Namespace, req.Name))

	log.V(3).Info("Fetching object")
	image := &anzalabsdevv1alpha1.Image{}
	if err := r.Get(ctx, req.NamespacedName, image); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if image.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(image, imageFinalizer) {
			log.V(3).Info("Adding finalizer")
			controllerutil.AddFinalizer(image, imageFinalizer)
			if err := r.Update(ctx, image); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(image, imageFinalizer) {
			log.V(1).Info("Deleting external resources")
			if err := r.deleteExternalResources(ctx, image); err != nil {
				return ctrl.Result{}, err
			}

			log.V(3).Info("Removing finalizer")
			controllerutil.RemoveFinalizer(image, imageFinalizer)
			if err := r.Update(ctx, image); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	log.V(1).Info("Building image", "format", image.Spec.Format)
	obj, err := r.buildImage(ctx, image)
	image.Status.Objects = obj
	if err != nil {
		if obj != nil {
			image.Status.Ready = false
			if statusErr := r.Status().Update(ctx, image); statusErr != nil {
				return ctrl.Result{}, errors.Join(err, statusErr)
			}
		}

		return ctrl.Result{}, err
	}

	log.V(1).Info("Image built and uploaded, updating status")
	image.Status.Ready = true
	if statusErr := r.Status().Update(ctx, image); statusErr != nil {
		return ctrl.Result{}, errors.Join(err, statusErr)
	}

	return ctrl.Result{}, nil
}

func (r *ImageReconciler) createClient(ctx context.Context, image *anzalabsdevv1alpha1.Image) (*s3.Client, error) {
	// Fetch the secret from image.Spec.BucketCredentials
	credentials := image.Spec.BucketCredentials
	if credentials.Name == "" {
		return nil, errors.New("bucket credentials are missing")
	}

	namespace := credentials.Namespace
	if namespace == "" {
		namespace = image.GetNamespace()
	}

	secret := &corev1.Secret{}
	if err := r.Get(
		ctx,
		types.NamespacedName{
			Namespace: namespace,
			Name:      credentials.Name,
		},
		secret,
	); err != nil {
		return nil, fmt.Errorf("failed to get secret %s: %w", credentials.Name, err)
	}

	// Unmarshal secret data into an s3.Config
	var config s3.Config
	if err := json.Unmarshal(secret.Data["BucketInfo.json"], &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bucket credentials: %w", err)
	}

	// Create an S3 client
	client, err := s3.New(config, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	return client, nil
}

func (r *ImageReconciler) deleteExternalResources(ctx context.Context, image *anzalabsdevv1alpha1.Image) error {
	log := log.FromContext(ctx, "image", klog.KObj(image))

	s3cli, err := r.createClient(ctx, image)
	if err != nil {
		return fmt.Errorf("failed to create S3 client: %w", err)
	}

	for objectKey := range image.Status.Objects {
		ok, err := s3cli.Stat(ctx, objectKey)
		if err != nil {
			return fmt.Errorf("failed to stat object: %w", err)
		} else if !ok {
			log.V(3).Info("Skipping, object does not exist", "object.key", objectKey)
			continue
		}

		if err := s3cli.Delete(ctx, objectKey); err != nil {
			return fmt.Errorf("failed to delete object: %w", err)
		}
	}

	return nil
}

func (r *ImageReconciler) buildImage(ctx context.Context, image *anzalabsdevv1alpha1.Image) (map[string]string, error) {
	log := log.FromContext(ctx, "image", klog.KObj(image))

	s3cli, err := r.createClient(ctx, image)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	out, err := r.Builder.Build(ctx, image.Spec.Format, image.Spec.Configuration)
	if err != nil {
		return nil, fmt.Errorf("failed to build image: %w", err)
	}

	objects := make(map[string]string)
	for _, out := range out {
		log.V(3).Info("Uploading object", "object.name", out.Name)

		f, err := os.Open(out.Path)
		if err != nil {
			return objects, fmt.Errorf("failed to open artifact %s: %w", out.Name, err)
		}

		objectKey := naming.Key(image, out.Name)
		if err := s3cli.Put(ctx, objectKey, f, out.Size); err != nil {
			return objects, fmt.Errorf("failed to upload image: %w", err)
		}
		objects[objectKey] = ""

		url, err := s3cli.GetURL(ctx, objectKey)
		if err != nil {
			return objects, fmt.Errorf("failed to generate url: %w", err)
		}
		objects[objectKey] = url
	}

	return objects, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ImageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&anzalabsdevv1alpha1.Image{}).
		Complete(r)
}
