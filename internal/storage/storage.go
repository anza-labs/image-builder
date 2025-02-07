// Copyright 2024-2025 anza-labs contributors.
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

	"github.com/anza-labs/image-builder/internal/storage/azure"
	"github.com/anza-labs/image-builder/internal/storage/s3"
)

var ErrInvalidConfig = errors.New("invalid configuration")

type Config struct {
	Spec Spec `json:"spec"`
}

type Spec struct {
	BucketName         string             `json:"bucketName"`
	AuthenticationType string             `json:"authenticationType"`
	Protocols          []string           `json:"protocols"`
	SecretS3           *s3.SecretS3       `json:"secretS3,omitempty"`
	SecretAzure        *azure.SecretAzure `json:"secretAzure,omitempty"`
}

type Storage interface {
	Delete(ctx context.Context, key string) error
	Get(ctx context.Context, key string, wr io.Writer) error
	GetURL(ctx context.Context, key string) (string, error)
	Put(ctx context.Context, key string, data io.Reader, size int64) error
	Stat(ctx context.Context, key string) (bool, error)
}

func New(config Config, ssl bool) (Storage, error) {
	// default to S3
	if slices.ContainsFunc(config.Spec.Protocols, func(s string) bool { return strings.EqualFold(s, "s3") }) {
		if !strings.EqualFold(config.Spec.AuthenticationType, "key") {
			return nil, fmt.Errorf("%w: invalid authentication type for s3", ErrInvalidConfig)
		}

		s3secret := config.Spec.SecretS3
		if s3secret == nil {
			return nil, fmt.Errorf("%w: s3 secret missing", ErrInvalidConfig)
		}

		return s3.New(config.Spec.BucketName, *s3secret, ssl)
	}

	// optionally Azure Blob
	if slices.ContainsFunc(config.Spec.Protocols, func(s string) bool { return strings.EqualFold(s, "azure") }) {
		if !strings.EqualFold(config.Spec.AuthenticationType, "key") {
			return nil, fmt.Errorf("%w: invalid authentication type for azure", ErrInvalidConfig)
		}

		azureSecret := config.Spec.SecretAzure
		if azureSecret == nil {
			return nil, fmt.Errorf("%w: azure secret missing", ErrInvalidConfig)
		}

		return azure.New(config.Spec.BucketName, *azureSecret)
	}

	return nil, fmt.Errorf("%w: invalid protocol (%v)", ErrInvalidConfig, config.Spec.Protocols)
}
