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

	"github.com/anza-labs/image-builder/internal/builder/linuxkit"
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
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

type options struct {
	Format             string
	ConfigPath         string
	OutputName         string
	StorageCredentials string
	K8sNamespace       string
	K8sJobName         string
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	ctrl.SetLogger(klog.NewKlogr())

	if err := run(signals.SetupSignalHandler(), options{
		Format:             os.Getenv("LINUXKIT_FORMAT"),
		ConfigPath:         os.Getenv("LINUXKIT_CONFIG"),
		StorageCredentials: os.Getenv("STORAGE_CREDENTIALS"),
		OutputName:         os.Getenv("K8S_SECRET_NAME"),
		K8sNamespace:       os.Getenv("K8S_NAMESPACE"),
		K8sJobName:         os.Getenv("K8S_JOB_NAME"),
	}); err != nil {
		klog.V(0).ErrorS(err, "Critical error while running")
		os.Exit(1)
	}
}

func run(ctx context.Context, opts options) error {
	log := log.FromContext(ctx)

	log.V(1).Info("Starting run", "options", opts)

	log.V(1).Info("Opening storage credentials file", "path", opts.StorageCredentials)
	f, err := os.Open(opts.StorageCredentials)
	if err != nil {
		return fmt.Errorf("failed to open BucketInfo.json: %w", err)
	}
	defer f.Close() //nolint:errcheck // best effort call

	var cfg storage.Config
	log.V(1).Info("Decoding storage credentials")
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return fmt.Errorf("failed to decode bucket credentials: %w", err)
	}

	log.V(1).Info("Creating Kubernetes client")
	cli, err := client.New(config.GetConfigOrDie(), client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	log.V(1).Info("Initializing storage")
	stor, err := storage.New(cfg, true)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	log.V(1).Info("Initializing builder")
	bld, err := linuxkit.New()
	if err != nil {
		return fmt.Errorf("failed to initialize builder: %w", err)
	}

	log.V(1).Info("Building images", "format", opts.Format, "configPath", opts.ConfigPath)
	out, err := bld.Build(ctx, opts.Format, opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to build images: %w", err)
	}

	outputs := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: opts.K8sNamespace,
			Name:      opts.OutputName,
		},
		Data: map[string][]byte{},
	}

	log.V(1).Info("Processing output objects", "objects", out)
	for _, o := range out {
		log.V(1).Info("Processing output file", "path", o.Path)
		f, err := os.Open(o.Path)
		if err != nil {
			return fmt.Errorf("failed to open file at path %s: %w", o.Path, err)
		}
		defer f.Close() //nolint:errcheck // best effort call

		objectKey := naming.Key(opts.K8sNamespace, opts.K8sJobName, opts.Format, o.Name)
		log.V(1).Info("Uploading image to storage", "key", objectKey)
		if err := stor.Put(ctx, objectKey, f, o.Size); err != nil {
			return fmt.Errorf("failed to upload image to storage with key %s: %w", objectKey, err)
		}

		log.V(1).Info("Generating URL for object", "key", objectKey)
		url, err := stor.GetURL(ctx, objectKey)
		if err != nil {
			return fmt.Errorf("failed to generate URL for object key %s: %w", objectKey, err)
		}

		if outputs.Data == nil {
			log.V(3).Info("Data map was empty, initializing")
			outputs.Data = map[string][]byte{}
		}

		log.V(6).Info("New data added to secret", "key", objectKey, "url", url)
		outputs.Data[naming.DNSName(o.Name)] = []byte(fmt.Sprintf("%s = %s", objectKey, url))
	}

	log.V(1).Info("Creating or updating Kubernetes secret", "secret", klog.KObj(outputs))
	err = cli.Create(ctx, outputs)
	if err != nil {
		return fmt.Errorf("failed to create or update Kubernetes secret: %w", err)
	}

	log.V(1).Info("Run completed successfully")
	return nil
}
