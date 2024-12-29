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

package builder

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Builder struct {
	linuxkit string
}

type Output struct {
	Size int64
	Path string
	Name string
}

func New() (*Builder, error) {
	linuxkit, err := exec.LookPath("linuxkit")
	if err != nil {
		if !errors.Is(err, exec.ErrNotFound) {
			return nil, fmt.Errorf("unable to look path: %w", err)
		}
		linuxkit = "/linuxkit"
	}

	return &Builder{
		linuxkit: linuxkit,
	}, nil
}

func FilePathWalkDir(root string) ([]Output, error) {
	var files []Output
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, Output{
				Path: path,
				Name: info.Name(),
				Size: info.Size(),
			})
		}
		return nil
	})
	return files, err
}

func (b *Builder) Build(ctx context.Context, format string, configPath string) ([]Output, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare output dir: %w", err)
	}

	cmd := exec.CommandContext(
		ctx,
		b.linuxkit,
		"build",
		"--format", format,
		"--dir", dir,
		configPath,
	)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w: %s", err, stderr.String())
	}

	outputs, err := FilePathWalkDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read contents of %s: %w", dir, err)
	}

	return outputs, nil
}
