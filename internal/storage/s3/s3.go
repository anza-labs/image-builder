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
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/anza-labs/image-builder/internal/util"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Client struct {
	s3cli      *minio.Client
	bucketName string
	expiry     time.Duration
}

type SecretS3 struct {
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	AccessKeyID     string `json:"accessKeyID"`
	AccessSecretKey string `json:"accessSecretKey"`
}

func New(bucketName string, s3secret SecretS3, ssl bool) (*Client, error) {
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
		bucketName: bucketName,
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

func (c *Client) Get(ctx context.Context, key string, wr io.Writer) error {
	pwr := &util.ProgressWriter{
		Underlying: wr,
		Log:        log.FromContext(ctx).WithName("ProgressWriter"),
	}

	obj, err := c.s3cli.GetObject(ctx, c.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	_, err = io.Copy(pwr, obj)
	return err
}

func (c *Client) Put(ctx context.Context, key string, data io.Reader, size int64) error {
	r := &util.ProgressReader{
		Underlying: data,
		TotalSize:  size,
		Log:        log.FromContext(ctx, "size.total", size).WithName("ProgressReader"),
	}

	_, err := c.s3cli.PutObject(ctx, c.bucketName, key, r, size, minio.PutObjectOptions{})

	return err
}

func (c *Client) GetURL(ctx context.Context, key string) (string, error) {
	url, err := c.s3cli.PresignedGetObject(ctx, c.bucketName, key, c.expiry, nil)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}
