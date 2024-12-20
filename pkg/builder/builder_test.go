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

package builder

import (
	"context"
	"os"
	"path"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed test/simple.yaml
var simple string

func TestBuild(t *testing.T) {
	b, err := New()
	require.NoError(t, err)

	out, err := b.Build(context.Background(), "kernel+initrd", simple)
	assert.NoError(t, err)
	assert.NotEmpty(t, out)

	defer func() {
		for _, o := range out {
			err := os.RemoveAll(o.Path)
			assert.NoError(t, err)
		}
		if len(out) != 0 {
			err = os.RemoveAll(path.Dir(out[0].Path))
			assert.NoError(t, err)
		}
	}()
}