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

	anzalabsdevv1alpha1 "github.com/anza-labs/image-builder/api/v1alpha1"
	"github.com/anza-labs/image-builder/version"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
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
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=image-builder.anza-labs.dev,resources=images,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=image-builder.anza-labs.dev,resources=images/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=image-builder.anza-labs.dev,resources=images/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=secrets/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=serviceaccounts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=serviceaccounts/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=configmaps/finalizers,verbs=update
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=jobs/finalizers,verbs=update
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles/finalizers,verbs=update
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ImageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx, "image", klog.KRef(req.Namespace, req.Name))

	log.V(3).Info("Fetching Image object")
	image := &anzalabsdevv1alpha1.Image{}
	if err := r.Get(ctx, req.NamespacedName, image); err != nil {
		log.V(0).Error(err, "Failed to fetch Image object")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Handle finalizer logic
	if image.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(image, imageFinalizer) {
			log.V(3).Info("Adding finalizer")
			controllerutil.AddFinalizer(image, imageFinalizer)
			if err := r.Update(ctx, image); err != nil {
				log.V(0).Error(err, "Failed to update Image object with finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(image, imageFinalizer) {
			// Perform cleanup
			log.V(3).Info("Performing cleanup and removing finalizer")
			if err := r.cleanupResources(ctx, image); err != nil {
				log.V(0).Error(err, "Failed to clean up resources")
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(image, imageFinalizer)
			if err := r.Update(ctx, image); err != nil {
				log.V(0).Error(err, "Failed to update Image object during finalizer removal")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	job := Job(image, version.Version)
	if err := r.ensureResources(ctx, image,
		ServiceAccount(image),
		Role(image),
		RoleBinding(image),
		ConfigMap(image),
		job,
	); err != nil {
		log.V(0).Error(err, "Failed to ensure resources")
		return ctrl.Result{}, err
	}

	// Update status based on Job completion
	log.V(3).Info("Checking Job completion")
	jobStatus := &batchv1.Job{}
	if err := r.Get(ctx, client.ObjectKeyFromObject(job), jobStatus); err != nil {
		log.V(0).Error(err, "Failed to fetch Job status")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if jobStatus.Status.Succeeded > 0 {
		log.V(3).Info("Job completed successfully")
		image.Status.Ready = true
		if err := r.Status().Update(ctx, image); err != nil {
			log.V(0).Error(err, "Failed to update Image status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// ensureResource ensures that a resource is created or updated.
func (r *ImageReconciler) ensureResources(ctx context.Context, owner client.Object, objs ...client.Object) error {
	log := log.FromContext(ctx, "owner", klog.KObj(owner))

	for _, resource := range objs {
		log.V(3).Info("Ensuring object exists",
			"name", resource.GetName(),
			"kind", resource.GetObjectKind().GroupVersionKind().Kind)

		_, err := controllerutil.CreateOrUpdate(ctx, r.Client, resource, func() error {
			return ctrl.SetControllerReference(owner, resource, r.Scheme)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// cleanupResources removes resources owned by the Image.
func (r *ImageReconciler) cleanupResources(ctx context.Context, image *anzalabsdevv1alpha1.Image) error {
	log := log.FromContext(ctx)
	log.V(3).Info("Cleaning up resources", "image", klog.KRef(image.Namespace, image.Name))

	// Define a list of owned resources to delete
	resourceTypes := []client.ObjectList{
		&batchv1.JobList{},
		&corev1.ConfigMapList{},
		&corev1.SecretList{},
		&corev1.ServiceAccountList{},
		&rbacv1.RoleList{},
		&rbacv1.RoleBindingList{},
	}

	for _, resourceType := range resourceTypes {
		list := resourceType.DeepCopyObject().(client.ObjectList)
		err := r.List(ctx, list, client.InNamespace(image.Namespace), client.MatchingFields{
			"metadata.ownerReferences": string(image.UID),
		})
		if err != nil {
			return err
		}

		// Iterate over resources and delete them
		items, err := meta.ExtractList(list)
		if err != nil {
			return err
		}
		for _, item := range items {
			resource := item.(client.Object)
			log.V(3).Info("Deleting resource",
				"name", resource.GetName(),
				"kind", resource.GetObjectKind().GroupVersionKind().Kind)
			if err := r.Delete(ctx, resource); err != nil {
				return err
			}
		}
	}

	log.V(3).Info("Cleanup complete")
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ImageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&anzalabsdevv1alpha1.Image{}).
		Owns(&batchv1.Job{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&rbacv1.Role{}).
		Owns(&rbacv1.RoleBinding{}).
		Complete(r)
}
