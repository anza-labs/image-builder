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

package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/anza-labs/image-builder/internal/storage/s3"
)

var ErrInvalidConfig = errors.New("invalid configuration")

type Config struct {
	Spec Spec `json:"spec"`
}

type Spec struct {
	BucketName         string       `json:"bucketName"`
	AuthenticationType string       `json:"authenticationType"`
	Protocols          []string     `json:"protocols"`
	SecretS3           *s3.SecretS3 `json:"secretS3,omitempty"`
}

type Storage interface {
	Delete(ctx context.Context, key string) error
	Get(ctx context.Context, key string, wr io.Writer) error
	GetURL(ctx context.Context, key string) (string, error)
	Put(ctx context.Context, key string, data io.Reader, size int64) error
	Stat(ctx context.Context, key string) (bool, error)
}

func New(config Config, ssl bool) (Storage, error) {
	if !slices.ContainsFunc(config.Spec.Protocols, func(s string) bool { return strings.EqualFold(s, "s3") }) {
		return nil, fmt.Errorf("%w: invalid protocol", ErrInvalidConfig)
	}

	if !strings.EqualFold(config.Spec.AuthenticationType, "key") {
		return nil, fmt.Errorf("%w: invalid authentication type", ErrInvalidConfig)
	}

	s3secret := config.Spec.SecretS3
	if s3secret == nil {
		return nil, fmt.Errorf("%w: s3 secret missing", ErrInvalidConfig)
	}

	return s3.New(config.Spec.BucketName, *config.Spec.SecretS3, ssl)
}
