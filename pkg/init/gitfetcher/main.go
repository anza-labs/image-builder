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
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/anza-labs/image-builder/internal/fetcherconfig"
	"github.com/anza-labs/image-builder/internal/git"
	"github.com/anza-labs/image-builder/internal/util"

	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

const (
	defaultUsername string = "gitfetcher"
)

var (
	ErrEmptyPasswordPath = errors.New("empty password file path")
)

type options struct {
	Config string
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	ctrl.SetLogger(klog.NewKlogr())

	if err := run(signals.SetupSignalHandler(), options{
		Config: os.Getenv("FETCHER_CONFIG"),
	}); err != nil {
		klog.V(0).ErrorS(err, "Critical error while running")
		os.Exit(1)
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
		if fetcher.GitFetcher == nil {
			log.V(4).Info("Ignoring fetcher config, not an GitFetcher")
			continue
		}

		if err := runFetcher(ctx, fetcher.GitFetcher); err != nil {
			log.V(1).Error(err, "New error occurred while running fetcher", "mount_point", fetcher.GitFetcher)
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return fmt.Errorf("one or more errors occurred: %w", err)
	}

	log.V(1).Info("Run completed successfully")
	return nil
}

func runFetcher(ctx context.Context, cfg *fetcherconfig.GitFetcher) error {
	log := log.FromContext(ctx)

	c, err := newClient(cfg.CredentialsPath)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	log.V(1).Info("Cloning repository", "repo", cfg.Repository, "ref", cfg.Ref)

	return c.Clone(ctx, cfg.Repository, cfg.Ref, cfg.MountPoint)
}

func newClient(credentialsPath string) (*git.Client, error) {
	entries, err := os.ReadDir(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config dir: %w", err)
	}

	var opts []git.Option

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		completePath := filepath.Join(credentialsPath, e.Name())

		switch e.Name() {
		case "password":
			usernamePath := filepath.Join(credentialsPath, "username")
			if _, err := os.Stat(usernamePath); err != nil {
				usernamePath = ""
			}
			auth, err := usernameAndPassword(usernamePath, completePath)
			if err != nil {
				return nil, fmt.Errorf("failed to create username/password auth: %w", err)
			}
			opts = append(opts, git.WithAuth(auth))

		case "gitconfig":
			cfg, err := gitConfig(completePath)
			if err != nil {
				return nil, fmt.Errorf("failed to load gitconfig: %w", err)
			}
			opts = append(opts, git.WithGitConfig(cfg))

		case "ssh-privatekey":
			sshAuth, err := sshPrivateKey(completePath)
			if err != nil {
				return nil, fmt.Errorf("failed to load SSH private key: %w", err)
			}
			opts = append(opts, git.WithAuth(sshAuth))
		}
	}

	client, err := git.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to create git client: %w", err)
	}

	return client, nil
}

func gitConfig(gitconfigFile string) (*config.Config, error) {
	b, err := util.ReadFile(gitconfigFile)
	if err != nil {
		return nil, err
	}

	cfg := config.NewConfig()
	if err := cfg.Unmarshal(b); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return cfg, nil
}

func usernameAndPassword(usernameFile, passwordFile string) (http.AuthMethod, error) {
	var username string
	var err error
	if usernameFile == "" {
		username = defaultUsername
	} else {
		b, err := util.ReadFile(usernameFile)
		if err != nil {
			return nil, err
		}
		username = string(b)
	}

	if passwordFile == "" {
		return nil, ErrEmptyPasswordPath
	}

	b, err := util.ReadFile(passwordFile)
	if err != nil {
		return nil, err
	}

	return &http.BasicAuth{
		Username: username,
		Password: string(b),
	}, nil
}

func sshPrivateKey(pemFile string) (ssh.AuthMethod, error) {
	b, err := util.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}

	pk, err := ssh.NewPublicKeys(defaultUsername, b, "")
	if err != nil {
		return nil, fmt.Errorf("unable to create public keys: %w", err)
	}

	return pk, nil
}
