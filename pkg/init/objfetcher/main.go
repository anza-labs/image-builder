// Copyright 2025 anza-labs contributors.
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
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anza-labs/image-builder/internal/fetcherconfig"
	"github.com/anza-labs/image-builder/internal/storage"
	"github.com/anza-labs/image-builder/internal/util"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

type options struct {
	K8sNamespace string
	K8sName      string
	Config       string
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	ctrl.SetLogger(klog.NewKlogr())

	if err := run(signals.SetupSignalHandler(), options{
		K8sNamespace: os.Getenv("K8S_NAMESPACE"),
		K8sName:      os.Getenv("K8S_NAME"),
		Config:       os.Getenv("FETCHER_CONFIG"),
	}); err != nil {
		klog.V(0).ErrorS(err, "Critical error while running")
	}
}

func run(ctx context.Context, opts options) error {
	log := log.FromContext(ctx)

	log.V(1).Info("Starting run", "options", opts)

	cfg, err := fetcherconfig.Load(opts.Config)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	var errs error
	for _, fetcher := range cfg.Fetchers {
		if fetcher.ObjFetcher == nil {
			continue
		}

		if err := runFetcher(ctx, fetcher.ObjFetcher); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return err // TODO
	}

	log.V(1).Info("Run completed successfully")
	return nil
}

func runFetcher(ctx context.Context, cfg *fetcherconfig.ObjFetcher) error {
	if cfg.KeysPath != "" {
		if err := loadKeys(cfg); err != nil {
			return fmt.Errorf("failed to load keys: %w", err)
		}
	}

	c, err := newClient(cfg.CredentialsPath)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for key, file := range cfg.Keys {
		err := saveObject(ctx, c, key, file)
		if err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	}

	return nil
}

func saveObject(ctx context.Context, client storage.Storage, key string, file fetcherconfig.File) error {
	f, err := os.OpenFile(file.Path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(file.Mode))
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", file.Path, err)
	}
	defer f.Close()

	if err := client.Get(ctx, key, f); err != nil {
		return fmt.Errorf("failed to fetch object with key %s: %w", key, err)
	}

	return nil
}

func newClient(credentialsPath string) (storage.Storage, error) {
	b, err := util.ReadFile(credentialsPath)
	if err != nil {
		return nil, err
	}

	var cfg storage.Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode bucket credentials: %w", err)
	}

	return storage.New(cfg, true)
}

func loadKeys(cfg *fetcherconfig.ObjFetcher) error {
	entries, err := os.ReadDir(cfg.KeysPath)
	if err != nil {
		return fmt.Errorf("failed to read keys directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		completePath := filepath.Join(cfg.KeysPath, entry.Name())
		if entry.Name() == "" || entry.Name() == "." || entry.Name() == ".." {
			return fmt.Errorf("invalid file name: %s", entry.Name())
		}

		data, err := os.ReadFile(completePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", completePath, err)
		}

		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			if key == "" {
				return fmt.Errorf("empty key found in file %s", entry.Name())
			}

			fp := filepath.Base(key)
			if fp == "" || fp == "." || fp == ".." {
				return fmt.Errorf("invalid key as file name: %s", key)
			}

			cfg.Keys[key] = fetcherconfig.File{
				Mode: 0o755,
				Path: fp,
			}
		}
	}

	return nil
}
