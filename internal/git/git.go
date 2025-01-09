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

package git

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

type Client struct {
	auth   transport.AuthMethod
	config *config.Config
}

func New(opts ...Option) (*Client, error) {
	cli := &Client{}

	var errs error
	for _, opt := range opts {
		if err := opt.apply(cli); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return cli, errs
}

type Option struct {
	apply func(*Client) error
}

func WithAuth(auth transport.AuthMethod) Option {
	return Option{
		apply: func(c *Client) error {
			if c.auth != nil {
				return errors.New("cannot set more than one auth method")
			}

			if auth == nil {
				return errors.New("auth cannot be nil")
			}

			c.auth = auth
			return nil
		},
	}
}

func WithGitConfig(gitconfig *config.Config) Option {
	return Option{
		apply: func(c *Client) error {
			if gitconfig == nil {
				return errors.New("gitconfig cannot be nil")
			}

			c.config = gitconfig
			return nil
		},
	}
}

func (c *Client) Clone(ctx context.Context, url, ref, path string) error {
	wt := osfs.New(path)
	dot, err := wt.Chroot(git.GitDirName)
	if err != nil {
		return fmt.Errorf("failed create git worktree dir: %w", err)
	}

	s := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
	if err := s.SetConfig(c.config); err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}

	_, err = git.CloneContext(ctx, s, wt, &git.CloneOptions{
		URL:               url,
		ReferenceName:     plumbing.ReferenceName(ref),
		SingleBranch:      true,
		Depth:             1,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              c.auth,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository %q at %q to %q: %w", url, ref, path, err)
	}

	return nil
}
