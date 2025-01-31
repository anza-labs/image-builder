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
	"testing"

	"github.com/distribution/reference"
	"github.com/stretchr/testify/assert"

	imagebuilderv1alpha2 "github.com/anza-labs/image-builder/api/v1alpha2"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConvertFromTo(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		src             *imagebuilderv1alpha2.Image
		expected        *Image
		expectedToErr   error
		expectedFromErr error
	}{
		"basic": {
			src: &imagebuilderv1alpha2.Image{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fromv1",
					Namespace: "test",
				},
				Spec: imagebuilderv1alpha2.ImageSpec{
					Builder: imagebuilderv1alpha2.Container{
						Image:     "ghcr.io/anza-labs/image-builder-linuxkit:v0.1.0",
						Verbosity: 4,
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								"cpu": resource.MustParse("1"),
							},
						},
					},
					GitFetcher: imagebuilderv1alpha2.Container{
						Image:     "ghcr.io/anza-labs/image-builder-init-gitfetcher:v0.1.0",
						Verbosity: 4,
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								"cpu": resource.MustParse("1"),
							},
						},
					},
					ObjFetcher: imagebuilderv1alpha2.Container{
						Image:     "ghcr.io/anza-labs/image-builder-init-objfetcher:v0.1.0",
						Verbosity: 4,
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								"cpu": resource.MustParse("1"),
							},
						},
					},
					Format:        "iso-efi",
					Configuration: "test",
					Result: corev1.LocalObjectReference{
						Name: "result",
					},
					BucketCredentials: corev1.LocalObjectReference{
						Name: "credentials",
					},
				},
				Status: imagebuilderv1alpha2.ImageStatus{
					Ready: true,
				},
			},
			expected: &Image{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fromv1",
					Namespace: "test",
				},
				Spec: ImageSpec{
					BuilderImage:     "ghcr.io/anza-labs/image-builder-linuxkit:v0.1.0",
					BuilderVerbosity: 4,
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							"cpu": resource.MustParse("1"),
						},
					},
					Format:        "iso-efi",
					Configuration: "test",
					Result: corev1.LocalObjectReference{
						Name: "result",
					},
					BucketCredentials: corev1.LocalObjectReference{
						Name: "credentials",
					},
				},
				Status: ImageStatus{
					Ready: true,
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var err error

			actual := &Image{}
			err = actual.ConvertFrom(tc.src)
			assert.ErrorIs(t, err, tc.expectedFromErr)
			assert.Equal(t, tc.expected, actual)

			actualHub := &imagebuilderv1alpha2.Image{}
			err = actual.ConvertTo(actualHub)
			assert.ErrorIs(t, err, tc.expectedToErr)
			assert.Equal(t, tc.src, actualHub)
		})
	}
}

func TestImageFrom(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		v1alpha1Image string
		targetImage   string
		expectedOut   string
		expectedErr   error
	}{
		"empty name and builderImage": {
			v1alpha1Image: "",
			targetImage:   "",
			expectedOut:   "",
		},
		"valid image": {
			v1alpha1Image: "image-builder-linuxkit:latest",
			targetImage:   "image-builder-init-objfetcher",
			expectedOut:   "docker.io/library/image-builder-init-objfetcher:latest",
		},
		"valid image with tag": {
			v1alpha1Image: "ghcr.io/anza-labs/image-builder-linuxkit:latest",
			targetImage:   "image-builder-init-objfetcher",
			expectedOut:   "ghcr.io/anza-labs/image-builder-init-objfetcher:latest",
		},
		"valid image without registry": {
			v1alpha1Image: "anza-labs/image-builder-linuxkit:latest",
			targetImage:   "image-builder-init-objfetcher",
			expectedOut:   "docker.io/anza-labs/image-builder-init-objfetcher:latest",
		},
		"valid image without local registry": {
			v1alpha1Image: "localhost:5005/anza-labs/image-builder-linuxkit:latest",
			targetImage:   "image-builder-init-objfetcher",
			expectedOut:   "localhost:5005/anza-labs/image-builder-init-objfetcher:latest",
		},
		"valid image with digest": {
			v1alpha1Image: "ghcr.io/anza-labs/image-builder@sha256:008b026f11c0b5653d564d0c9877a116770f06dfbdb36ca75c46fd593d863cbc",
			targetImage:   "image-builder-init-objfetcher",
			expectedOut:   "ghcr.io/anza-labs/image-builder-init-objfetcher@sha256:008b026f11c0b5653d564d0c9877a116770f06dfbdb36ca75c46fd593d863cbc",
		},
		"valid image with tag and digest": {
			v1alpha1Image: "ghcr.io/anza-labs/image-builder-linuxkit:v0.1.2@sha256:008b026f11c0b5653d564d0c9877a116770f06dfbdb36ca75c46fd593d863cbc",
			targetImage:   "image-builder-init-objfetcher",
			expectedOut:   "ghcr.io/anza-labs/image-builder-init-objfetcher:v0.1.2@sha256:008b026f11c0b5653d564d0c9877a116770f06dfbdb36ca75c46fd593d863cbc",
		},
		"invalid builderImage": {
			v1alpha1Image: "!!!invalid_image",
			targetImage:   "image-builder-init-objfetcher",
			expectedErr:   reference.ErrReferenceInvalidFormat,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual, err := imageFrom(tc.v1alpha1Image, tc.targetImage)
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedOut, actual)
		})
	}
}
