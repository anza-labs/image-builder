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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ImageSpec defines the desired state of Image.
type ImageSpec struct {
	// BuilderImage indicates the container image to use for the Builder job.
	// +optional
	BuilderImage string `json:"builderImage,omitempty"`

	// Resources describe the compute resource requirements.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Affinity specifies the scheduling constraints for Pods.
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// Format specifies the image format.
	// +kubebuilder:validation:Enum=aws;docker;dynamic-vhd;gcp;iso-bios;iso-efi;iso-efi-initrd;kernel+initrd;kernel+iso;kernel+squashfs;qcow2-bios;qcow2-efi;raw-bios;raw-efi;rpi3;tar;tar-kernel-initrd;vhd;vmdk
	// +required
	Format string `json:"format"`

	// Configuration is a YAML formatted Linuxkit config.
	// +required
	Configuration string `json:"configuration"`

	// Result is a local reference that lists downloadable objects, that are results of the image building.
	// Defaults to the Image.Metadata.Name.
	// +optional
	Result corev1.LocalObjectReference `json:"result"`

	// BucketCredentials is a reference to the credentials for S3, where the image will be stored.
	// +required
	BucketCredentials corev1.LocalObjectReference `json:"bucketCredentials"`
}

// ImageStatus defines the observed state of Image.
type ImageStatus struct {
	// Ready indicates whether the image is ready.
	// +optional
	Ready bool `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready"

// Image is the Schema for the images API.
type Image struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageSpec   `json:"spec,omitempty"`
	Status ImageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ImageList contains a list of Image.
type ImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Image `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Image{}, &ImageList{})
}
