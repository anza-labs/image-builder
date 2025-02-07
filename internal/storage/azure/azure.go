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

package azure

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

	"github.com/anza-labs/image-builder/internal/util"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Client struct {
	azCli         *azblob.Client
	containerName string
}

type SecretAzure struct {
	AccessToken     string    `json:"accessToken"`
	ExpiryTimestamp time.Time `json:"expiryTimeStamp"`
}

func New(containerName string, azureSecret SecretAzure) (*Client, error) {
	azCli, err := azblob.NewClientWithNoCredential(azureSecret.AccessToken, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %w", err)
	}

	return &Client{
		azCli:         azCli,
		containerName: containerName,
	}, nil
}

func (c *Client) Stat(ctx context.Context, blobName string) (bool, error) {
	pager := c.azCli.NewListBlobsFlatPager(c.containerName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return false, fmt.Errorf("unable to fetch next page of results: %w", err)
		}

		segment := page.ListBlobsFlatSegmentResponse.Segment
		if segment == nil {
			return false, fmt.Errorf("segment is missing")
		}

		for _, item := range segment.BlobItems {
			if item == nil || item.Name == nil {
				continue
			}
			if *item.Name == blobName {
				return true, nil
			}
		}
	}

	return false, nil
}

func (c *Client) Delete(ctx context.Context, blobName string) error {
	_, err := c.azCli.DeleteBlob(ctx, c.containerName, blobName, nil)
	return err
}

func (c *Client) Get(ctx context.Context, blobName string, wr io.Writer) error {
	stream, err := c.azCli.DownloadStream(ctx, c.containerName, blobName, nil)
	if err != nil {
		return fmt.Errorf("unable to get download stream: %w", err)
	}

	pwr := &util.ProgressWriter{
		Underlying: wr,
		Log:        log.FromContext(ctx).WithName("ProgressWriter"),
	}

	_, err = io.Copy(pwr, stream.Body)
	return err
}

func (c *Client) Put(ctx context.Context, blobName string, data io.Reader, size int64) error {
	r := &util.ProgressReader{
		Underlying: data,
		TotalSize:  size,
		Log:        log.FromContext(ctx, "size.total", size).WithName("ProgressReader"),
	}

	_, err := c.azCli.UploadStream(ctx, c.containerName, blobName, r, nil)
	return err
}

func (c *Client) GetURL(ctx context.Context, blobName string) (string, error) {
	return "", errors.ErrUnsupported
}
