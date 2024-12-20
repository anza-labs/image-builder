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
	"path"

	anzalabsdevv1alpha1 "github.com/anza-labs/image-builder/api/v1alpha1"
)

func Key(image *anzalabsdevv1alpha1.Image, key string) string {
	return path.Clean(path.Join(
		DNSName(image.GetNamespace()),
		DNSName(image.GetName()),
		DNSName(image.Spec.Format),
		DNSName(key),
	))
}
