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

package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/go-logr/logr"
)

const (
	defaultLimit = 4 * 1024 // 4KiB
)

func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open %q file: %w", path, err)
	}
	defer f.Close() //nolint:errcheck // best effort call

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, io.LimitReader(f, defaultLimit))
	if err != nil {
		return nil, fmt.Errorf("unable to read %q file: %w", path, err)
	}

	return buf.Bytes(), nil
}

type ProgressReader struct {
	Log         logr.Logger
	Underlying  io.Reader
	TotalSize   int64
	currentSize int64
}

func (r *ProgressReader) Read(p []byte) (int, error) {
	n, err := r.Underlying.Read(p)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			r.Log.V(5).Error(err, "Underlying read errored",
				"size.current", r.currentSize+int64(n))
		}
		return n, err
	}

	r.currentSize += int64(n)

	var percentage int64
	if r.TotalSize > 0 {
		percentage = r.currentSize * 100 / r.TotalSize
	}

	r.Log.V(5).Info("Read successful, progressing",
		"size.current", r.currentSize,
		"size.percentage", percentage)
	return n, err
}

type ProgressWriter struct {
	Log         logr.Logger
	Underlying  io.Writer
	TotalSize   int64
	currentSize int64
}

func (w *ProgressWriter) Write(p []byte) (int, error) {
	n, err := w.Underlying.Write(p)
	if err != nil {
		w.Log.V(5).Error(err, "Underlying write errored",
			"size.current", w.currentSize+int64(n))
		return n, err
	}

	w.currentSize += int64(n)

	var percentage int64
	if w.TotalSize > 0 {
		percentage = w.currentSize * 100 / w.TotalSize
	}

	w.Log.V(5).Info("Write successful, progressing",
		"size.current", w.currentSize,
		"size.percentage", percentage)
	return n, err
}
