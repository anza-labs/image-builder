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

package v1alpha1

import (
	"fmt"
	"path"
	"strings"

	imagebuilderv1alpha2 "github.com/anza-labs/image-builder/api/v1alpha2"
	"github.com/distribution/reference"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this Image to the Hub version.
func (src *Image) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*imagebuilderv1alpha2.Image)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Spec
	srcImage := src.Spec.BuilderImage
	verbosity := src.Spec.BuilderVerbosity
	resources := src.Spec.Resources

	builderImage, err := imageFrom(srcImage, "image-builder")
	if err != nil {
		return fmt.Errorf("container image conversion failed: %w", err)
	}
	dst.Spec.Builder = imagebuilderv1alpha2.Container{
		Image:     builderImage,
		Verbosity: verbosity,
		Resources: resources,
	}

	gitImage, err := imageFrom(srcImage, "image-builder-init-gitfetcher")
	if err != nil {
		return fmt.Errorf("container image conversion failed: %w", err)
	}

	dst.Spec.GitFetcher = imagebuilderv1alpha2.Container{
		Image:     gitImage,
		Verbosity: verbosity,
		Resources: resources,
	}

	objImage, err := imageFrom(srcImage, "image-builder-init-objfetcher")
	if err != nil {
		return fmt.Errorf("container image conversion failed: %w", err)
	}

	dst.Spec.ObjFetcher = imagebuilderv1alpha2.Container{
		Image:     objImage,
		Verbosity: verbosity,
		Resources: resources,
	}

	dst.Spec.Affinity = src.Spec.Affinity
	dst.Spec.BucketCredentials = src.Spec.BucketCredentials
	dst.Spec.Result = src.Spec.Result
	dst.Spec.Configuration = src.Spec.Configuration
	dst.Spec.Format = src.Spec.Format

	// Status
	dst.Status.Ready = src.Status.Ready

	return nil
}

// ConvertFrom converts from the Hub version to this version.
func (dst *Image) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*imagebuilderv1alpha2.Image)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Spec
	dst.Spec.BuilderImage = src.Spec.Builder.Image
	dst.Spec.BuilderVerbosity = src.Spec.Builder.Verbosity
	dst.Spec.Resources = src.Spec.Builder.Resources

	dst.Spec.Affinity = src.Spec.Affinity
	dst.Spec.BucketCredentials = src.Spec.BucketCredentials
	dst.Spec.Result = src.Spec.Result
	dst.Spec.Configuration = src.Spec.Configuration
	dst.Spec.Format = src.Spec.Format

	// Status
	dst.Status.Ready = src.Status.Ready

	return nil
}

// imageFrom combines a source image and a target image to generate a new image reference.
// It handles image names with or without registry, tags, and digests.
func imageFrom(v1alpha1Image, targetImage string) (string, error) {
	// If either the v1alpha1Image or targetImage is empty, return the targetImage as the result
	if v1alpha1Image == "" || targetImage == "" {
		return "", nil
	}

	// Parse the v1alpha1Image using Docker reference package
	ref, err := reference.ParseNormalizedNamed(v1alpha1Image)
	if err != nil {
		return "", err
	}

	// Extract the registry, repository, tag, and digest from the v1alpha1Image
	registry := reference.Domain(ref)
	baseImage := reference.Path(ref)

	if paths := strings.Split(baseImage, "/"); len(paths) > 1 {
		baseImage = strings.Join(paths[:len(paths)-1], "/")
	}

	// Combine the targetImage with the base image's registry/repository, tag, and digest
	var result strings.Builder

	result.WriteString(path.Join(registry, baseImage, targetImage))

	if tagged, ok := ref.(reference.Tagged); ok {
		result.WriteString(":" + tagged.Tag())
	}
	if digested, ok := ref.(reference.Digested); ok {
		result.WriteString("@" + digested.Digest().String())
	}

	// Return the resulting image
	return result.String(), nil
}
