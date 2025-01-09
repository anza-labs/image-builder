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

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anza-labs/image-builder/internal/fetcherconfig"
)

const (
	testDir = "test/keys"
)

func TestLoadKeys(t *testing.T) {
	t.Parallel()

	// Prepare
	cfg := &fetcherconfig.ObjFetcher{
		KeysPath: testDir,
	}
	expected := map[string]fetcherconfig.File{
		"key/of/obj1": {Path: "obj1", Mode: 0o755},
		"key/of/obj2": {Path: "obj2", Mode: 0o755},
	}

	// Test
	err := loadKeys(cfg)

	// Validate
	assert.NoError(t, err)
	assert.Equal(t, cfg.Keys, expected)
}
