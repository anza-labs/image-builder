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

package naming

import (
	"testing"

	"github.com/anza-labs/image-builder/api/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKey(t *testing.T) {
	tests := []struct {
		name           string
		format         string
		expectedOutput string
	}{
		{"AWS format", "aws", "test/test-image/aws/stdin"},
		{"Docker format", "docker", "test/test-image/docker/stdin"},
		{"Dynamic VHD format", "dynamic-vhd", "test/test-image/dynamic-vhd/stdin"},
		{"GCP format", "gcp", "test/test-image/gcp/stdin"},
		{"ISO BIOS format", "iso-bios", "test/test-image/iso-bios/stdin"},
		{"ISO EFI format", "iso-efi", "test/test-image/iso-efi/stdin"},
		{"ISO EFI Initrd format", "iso-efi-initrd", "test/test-image/iso-efi-initrd/stdin"},
		{"Kernel+Initrd format", "kernel+initrd", "test/test-image/kernel-initrd/stdin"},
		{"Kernel+ISO format", "kernel+iso", "test/test-image/kernel-iso/stdin"},
		{"Kernel+SquashFS format", "kernel+squashfs", "test/test-image/kernel-squashfs/stdin"},
		{"QCOW2 BIOS format", "qcow2-bios", "test/test-image/qcow2-bios/stdin"},
		{"QCOW2 EFI format", "qcow2-efi", "test/test-image/qcow2-efi/stdin"},
		{"RAW BIOS format", "raw-bios", "test/test-image/raw-bios/stdin"},
		{"RAW EFI format", "raw-efi", "test/test-image/raw-efi/stdin"},
		{"RPI3 format", "rpi3", "test/test-image/rpi3/stdin"},
		{"TAR format", "tar", "test/test-image/tar/stdin"},
		{"TAR Kernel+Initrd format", "tar-kernel-initrd", "test/test-image/tar-kernel-initrd/stdin"},
		{"VHD format", "vhd", "test/test-image/vhd/stdin"},
		{"VMDK format", "vmdk", "test/test-image/vmdk/stdin"},
		{"Empty format", "", "test/test-image/a/stdin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image := &v1alpha1.Image{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-image",
					Namespace: "test",
				},
				Spec: v1alpha1.ImageSpec{
					Format: tt.format,
				},
			}
			output := Key(image, "stdin")
			if output != tt.expectedOutput {
				t.Errorf("expected %s, got %s", tt.expectedOutput, output)
			}
		})
	}
}
