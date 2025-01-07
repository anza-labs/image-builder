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

package v1alpha2

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ImageSpec defines the desired state of an Image resource.
type ImageSpec struct {
	// Builder specifies the parameters for the main container configuration.
	// +optional
	Builder Container `json:"builder,omitempty"`

	// ObjFetcher specifies the parameters for the Object Fetcher init container configuration.
	// +optional
	ObjFetcher Container `json:"objFetcher,omitempty"`

	// GitFetcher specifies the parameters for the Git Fetcher init container configuration.
	// +optional
	GitFetcher Container `json:"gitFetcher,omitempty"`

	// Affinity specifies the scheduling constraints for Pods running the builder job.
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// Format specifies the output image format.
	// +kubebuilder:validation:Enum=aws;docker;dynamic-vhd;gcp;iso-bios;iso-efi;iso-efi-initrd;kernel+initrd;kernel+iso;kernel+squashfs;qcow2-bios;qcow2-efi;raw-bios;raw-efi;rpi3;tar;tar-kernel-initrd;vhd;vmdk
	// +required
	Format string `json:"format"`

	// Configuration is a YAML-formatted Linuxkit configuration.
	// +required
	Configuration string `json:"configuration"`

	// Result is a reference to the local object containing downloadable build results.
	// Defaults to the Image.Metadata.Name if not specified.
	// +optional
	Result corev1.LocalObjectReference `json:"result"`

	// BucketCredentials is a reference to the credentials used for storing the image in S3.
	// +required
	BucketCredentials corev1.LocalObjectReference `json:"bucketCredentials"`

	// AdditionalData specifies additional data sources required for building the image.
	// +optional
	AdditionalData []AdditionalData `json:"additionalData"`
}

// AdditionalData represents additional data sources for image building.
type AdditionalData struct {
	// Name specifies unique name for the additional data.
	// +required
	Name string `json:"name"`

	// VolumeMountPoint specifies the path where this data should be mounted.
	// +required
	VolumeMountPoint string `json:"volumeMountPoint"`

	// DataSource specifies the data source details.
	DataSource `json:",inline"`
}

// DataSource defines the available sources for additional data.
// Each data source is either used directly as a Volume for the image, or
// will be fetched into empty dir shared between init container and the builder.
type DataSource struct {
	// ConfigMap specifies a ConfigMap as a data source.
	// +optional
	ConfigMap *corev1.ConfigMapVolumeSource `json:"configMap,omitempty"`

	// Secret specifies a Secret as a data source.
	// +optional
	Secret *corev1.SecretVolumeSource `json:"secret,omitempty"`

	// Image specifies a container image as a data source.
	// +optional
	Image *corev1.ImageVolumeSource `json:"image,omitempty"`

	// Volume specifies a PersistentVolumeClaim as a data source.
	// +optional
	Volume *corev1.PersistentVolumeClaimVolumeSource `json:"volume,omitempty"`

	// Bucket specifies an S3 bucket as a data source.
	// +optional
	Bucket *BucketDataSource `json:"bucket,omitempty"`

	// GitRepository specifies a Git repository as a data source.
	// +optional
	GitRepository *GitRepository `json:"gitRepository,omitempty"`
}

// BucketDataSource represents an S3 bucket data source.
type BucketDataSource struct {
	// Credentials is a reference to the credentials for accessing the bucket.
	// +required
	Credentials *corev1.LocalObjectReference `json:"credentials"`

	// Items specifies specific items within the bucket to include.
	// +optional
	Items []corev1.KeyToPath `json:"items,omitempty"`

	// ItemsSecret specifies a Scret mapping item names to object storage keys.
	// Each value should either be a key of the object or follow the format "key = <Presigned URL>",
	// e.g.:
	//	item-1: "path/to/item-1 = <Presigned URL>"
	//	item-2: "path/to/item-2"
	// +optional
	ItemsSecret *corev1.LocalObjectReference `json:"itemsConfigMap,omitempty"`
}

// GitRepository represents a Git repository data source.
type GitRepository struct {
	// Repository specifies the URL of the Git repository.
	// +required
	Repository string `json:"repository"`

	// Credentials specifies the credentials for accessing the repository.
	// Secret must be one of the following types:
	// 	- "kubernetes.io/basic-auth" with "username" and "password" fields;
	// 	- "kubernetes.io/ssh-auth" with "ssh-privatekey" field;
	// 	- "Opaque" with "gitconfig" field.
	// +optional
	Credentials *corev1.LocalObjectReference `json:"credentials,omitempty"`
}

// ImageStatus defines the observed state of an Image resource.
type ImageStatus struct {
	// Ready indicates whether the image has been successfully built.
	// +optional
	Ready bool `json:"ready"`
}

type Container struct {
	// Image indicates the container image to use for the init container.
	// +optional
	Image string `json:"image,omitempty"`

	// Verbosity specifies the log verbosity level for the container.
	// +optional
	// +default=4
	// +kubebuilder:default=4
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=10
	Verbosity uint8 `json:"verbosity"`

	// Resources describe the compute resource requirements for the builder job.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready"

// Image represents the schema for the images API.
type Image struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageSpec   `json:"spec,omitempty"`
	Status ImageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ImageList contains a list of Image resources.
type ImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Image `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Image{}, &ImageList{})
}
