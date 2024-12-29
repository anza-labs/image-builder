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
)

func TestKey(t *testing.T) {
	tests := []struct {
		name           string
		format         string
		expectedOutput string
	}{
		{"AWS format", "aws", "test-namespace/test-image/aws/stdin"},
		{"Docker format", "docker", "test-namespace/test-image/docker/stdin"},
		{"Dynamic VHD format", "dynamic-vhd", "test-namespace/test-image/dynamic-vhd/stdin"},
		{"GCP format", "gcp", "test-namespace/test-image/gcp/stdin"},
		{"ISO BIOS format", "iso-bios", "test-namespace/test-image/iso-bios/stdin"},
		{"ISO EFI format", "iso-efi", "test-namespace/test-image/iso-efi/stdin"},
		{"ISO EFI Initrd format", "iso-efi-initrd", "test-namespace/test-image/iso-efi-initrd/stdin"},
		{"Kernel+Initrd format", "kernel+initrd", "test-namespace/test-image/kernel-initrd/stdin"},
		{"Kernel+ISO format", "kernel+iso", "test-namespace/test-image/kernel-iso/stdin"},
		{"Kernel+SquashFS format", "kernel+squashfs", "test-namespace/test-image/kernel-squashfs/stdin"},
		{"QCOW2 BIOS format", "qcow2-bios", "test-namespace/test-image/qcow2-bios/stdin"},
		{"QCOW2 EFI format", "qcow2-efi", "test-namespace/test-image/qcow2-efi/stdin"},
		{"RAW BIOS format", "raw-bios", "test-namespace/test-image/raw-bios/stdin"},
		{"RAW EFI format", "raw-efi", "test-namespace/test-image/raw-efi/stdin"},
		{"RPI3 format", "rpi3", "test-namespace/test-image/rpi3/stdin"},
		{"TAR format", "tar", "test-namespace/test-image/tar/stdin"},
		{"TAR Kernel+Initrd format", "tar-kernel-initrd", "test-namespace/test-image/tar-kernel-initrd/stdin"},
		{"VHD format", "vhd", "test-namespace/test-image/vhd/stdin"},
		{"VMDK format", "vmdk", "test-namespace/test-image/vmdk/stdin"},
		{"Empty format", "", "test-namespace/test-image/stdin"}, // invalid, but stays for test purposes
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := Key("test-namespace", "test-image", tt.format, "stdin")
			if output != tt.expectedOutput {
				t.Errorf("expected %s, got %s", tt.expectedOutput, output)
			}
		})
	}
}
