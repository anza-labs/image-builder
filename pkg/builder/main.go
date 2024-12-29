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

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	anzalabsdevv1alpha1 "github.com/anza-labs/image-builder/api/v1alpha1"
	"github.com/anza-labs/image-builder/internal/builder"
	"github.com/anza-labs/image-builder/internal/naming"
	"github.com/anza-labs/image-builder/internal/storage"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(anzalabsdevv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

type options struct {
	format             string
	configPath         string
	outputName         string
	storageCredentials string
	k8sNamespace       string
	k8sName            string
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	ctrl.SetLogger(klog.NewKlogr())

	if err := run(signals.SetupSignalHandler(), options{
		format:             os.Getenv("LINUXKIT_FORMAT"),
		configPath:         os.Getenv("LINUXKIT_CONFIG"),
		storageCredentials: os.Getenv("STORAGE_CREDENTIALS"),
		outputName:         os.Getenv("K8S_SECRET_NAME"),
		k8sNamespace:       os.Getenv("K8S_NAMESPACE"),
		k8sName:            os.Getenv("K8S_NAME"),
	}); err != nil {
		setupLog.V(0).Error(err, "Critical error while running")
	}
}

func run(ctx context.Context, opts options) error {
	setupLog.V(1).Info("Starting run", "options", opts)

	setupLog.V(1).Info("Opening storage credentials file", "path", opts.storageCredentials)
	f, err := os.Open(opts.storageCredentials)
	if err != nil {
		return fmt.Errorf("failed to open BucketInfo.json: %w", err)
	}
	defer f.Close()

	var cfg storage.Config
	setupLog.V(1).Info("Decoding storage credentials")
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return fmt.Errorf("failed to decode bucket credentials: %w", err)
	}

	setupLog.V(1).Info("Creating Kubernetes client")
	cli, err := client.New(config.GetConfigOrDie(), client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	setupLog.V(1).Info("Initializing storage")
	stor, err := storage.New(cfg, true)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	setupLog.V(1).Info("Initializing builder")
	bld, err := builder.New()
	if err != nil {
		return fmt.Errorf("failed to initialize builder: %w", err)
	}

	setupLog.V(1).Info("Building images", "format", opts.format, "configPath", opts.configPath)
	out, err := bld.Build(ctx, opts.format, opts.configPath)
	if err != nil {
		return fmt.Errorf("failed to build images: %w", err)
	}

	outputs := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: opts.k8sNamespace,
			Name:      opts.outputName,
		},
		StringData: make(map[string]string),
	}
	owner := &anzalabsdevv1alpha1.Image{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: opts.k8sNamespace,
			Name:      opts.k8sName,
		},
	}

	for _, o := range out {
		setupLog.V(1).Info("Processing output file", "path", o.Path)
		f, err := os.Open(o.Path)
		if err != nil {
			return fmt.Errorf("failed to open file at path %s: %w", o.Path, err)
		}
		defer f.Close()

		objectKey := naming.Key(opts.k8sNamespace, opts.k8sName, opts.format, o.Name)
		setupLog.V(1).Info("Uploading image to storage", "key", objectKey)
		if err := stor.Put(ctx, objectKey, f, o.Size); err != nil {
			return fmt.Errorf("failed to upload image to storage with key %s: %w", objectKey, err)
		}

		setupLog.V(1).Info("Generating URL for object", "key", objectKey)
		url, err := stor.GetURL(ctx, objectKey)
		if err != nil {
			return fmt.Errorf("failed to generate URL for object key %s: %w", objectKey, err)
		}

		outputs.StringData[objectKey] = url

		setupLog.V(1).Info("Creating or updating Kubernetes secret", "name", opts.outputName)
		_, err = controllerutil.CreateOrUpdate(ctx, cli, outputs, func() error {
			return ctrl.SetControllerReference(owner, outputs, scheme)
		})
		if err != nil {
			return fmt.Errorf("failed to create or update Kubernetes secret: %w", err)
		}
	}

	setupLog.V(1).Info("Run completed successfully")
	return nil
}
