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

package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	s3cli      *minio.Client
	bucketName string
	expiry     time.Duration
}

var ErrInvalidConfig = errors.New("invalid configuration")

type Config struct {
	Spec Spec `json:"spec"`
}

type Spec struct {
	BucketName         string    `json:"bucketName"`
	AuthenticationType string    `json:"authenticationType"`
	Protocols          []string  `json:"protocols"`
	SecretS3           *SecretS3 `json:"secretS3,omitempty"`
}

type SecretS3 struct {
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	AccessKeyID     string `json:"accessKeyID"`
	AccessSecretKey string `json:"accessSecretKey"`
}

func New(config Config, ssl bool) (*Client, error) {
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

	s3cli, err := minio.New(s3secret.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3secret.AccessKeyID, s3secret.AccessSecretKey, ""),
		Region: s3secret.Region,
		Secure: ssl,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %w", err)
	}

	return &Client{
		s3cli:      s3cli,
		bucketName: config.Spec.BucketName,
		expiry:     time.Hour * 24 * 5,
	}, nil
}

func (c *Client) Stat(ctx context.Context, key string) (bool, error) {
	_, err := c.s3cli.StatObject(ctx, c.bucketName, key, minio.StatObjectOptions{})
	if err != nil {
		merr := minio.ToErrorResponse(err)
		if merr.StatusCode == http.StatusNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	return c.s3cli.RemoveObject(ctx, c.bucketName, key, minio.RemoveObjectOptions{})
}

func (c *Client) Put(ctx context.Context, key string, data io.Reader, size int64) error {
	_, err := c.s3cli.PutObject(ctx, c.bucketName, key, data, size, minio.PutObjectOptions{})

	return err
}

func (c *Client) GetURL(ctx context.Context, key string) (string, error) {
	url, err := c.s3cli.PresignedGetObject(ctx, c.bucketName, key, c.expiry, nil)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}
