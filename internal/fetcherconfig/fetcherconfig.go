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

package fetcherconfig

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	return LoadFrom(f)
}

func LoadFrom(r io.Reader) (*Config, error) {
	cfg := &Config{}
	if err := json.NewDecoder(r).Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return cfg, nil
}

type Config struct {
	Fetchers []Fetcher `json:"fetchers"`
}

type Fetcher struct {
	GitFetcher *GitFetcher `json:"gitfetcher,omitempty"`
	ObjFetcher *ObjFetcher `json:"objfetcher,omitempty"`
}

type GitFetcher struct {
	MountPoint      string `json:"mountPoint"`
	CredentialsPath string `json:"credentialsPath"`
	Repository      string `json:"repository"`
	Ref             string `json:"ref"`
}

type ObjFetcher struct {
	MountPoint      string          `json:"mountPoint"`
	CredentialsPath string          `json:"credentialsPath"`
	KeysPath        string          `json:"keysPath,omitempty"`
	Keys            map[string]File `json:"keys"`
}

type File struct {
	Path string `json:"path"`
	Mode int32  `json:"mode"`
}
