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

	"github.com/go-logr/logr"
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

type progressReader struct {
	log         logr.Logger
	underlying  io.Reader
	totalSize   int64
	currentSize int64
}

func (r *progressReader) Read(p []byte) (int, error) {
	n, err := r.underlying.Read(p)
	if err != nil {
		r.log.V(5).Error(err, "Underlying read errored",
			"size.current", r.currentSize+int64(n))
		return n, err
	}

	r.currentSize += int64(n)

	var percentage int64
	if r.totalSize > 0 {
		percentage = r.currentSize * 100 / r.totalSize
	}

	r.log.V(5).Info("Read successful, progressing",
		"size.current", r.currentSize,
		"size.percentage", percentage)
	return n, err
}

func (c *Client) Put(ctx context.Context, key string, data io.Reader, size int64) error {
	r := &progressReader{
		underlying: data,
		totalSize:  size,
		log:        log.FromContext(ctx, "size.total", size).WithName("progressReader"),
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
