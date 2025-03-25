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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const KindMkosi = "Mkosi"

// MkosiSpec defines the desired state of Mkosi.
type MkosiSpec struct {
}

// MkosiStatus defines the observed state of Mkosi.
type MkosiStatus struct {
	// Ready indicates whether the image has been successfully built.
	// +optional
	Ready bool `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Mkosi is the Schema for the mkosis API.
type Mkosi struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MkosiSpec   `json:"spec,omitempty"`
	Status MkosiStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MkosiList contains a list of Mkosi.
type MkosiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Mkosi `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Mkosi{}, &MkosiList{})
}
