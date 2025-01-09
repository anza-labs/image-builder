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
	"os"
	"testing"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/stretchr/testify/assert"
)

const (
	testPrivateKey = "test/ssh/ed25519"
	testUsername   = "test/plain/username"
	testPassword   = "test/plain/password"
)

func TestSSHPrivateKey(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		param       string
		actual      func(assert.TestingT, interface{}, ...interface{}) bool
		expectedErr error
	}{
		"valid": {
			param: testPrivateKey,
		},
		"no file": {
			param:       "not-exist",
			actual:      assert.Empty,
			expectedErr: os.ErrNotExist,
		},
	} {
		t.Run(name, func(t *testing.T) {
			// Prepare
			t.Parallel()
			if tc.actual == nil {
				tc.actual = assert.NotNil
			}

			// Test
			actual, err := sshPrivateKey(tc.param)

			// Validate
			assert.ErrorIs(t, err, tc.expectedErr)
			tc.actual(t, actual)
		})
	}
}

func TestUsernameAndPassword(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		usernameFile string
		passwordFile string
		expected     http.AuthMethod
		expectedErr  error
	}{
		"valid username and password files": {
			usernameFile: testUsername,
			passwordFile: testPassword,
			expected: &http.BasicAuth{
				Username: "test-username",
				Password: "test-password",
			},
		},
		"empty username file path": {
			usernameFile: "",
			passwordFile: testPassword,
			expected: &http.BasicAuth{
				Username: defaultUsername,
				Password: "test-password",
			},
		},
		"empty password file path": {
			usernameFile: testUsername,
			passwordFile: "",
			expectedErr:  ErrEmptyPasswordPath,
		},
		"non-existent username file": {
			usernameFile: "nonexistent-username.txt",
			passwordFile: testPassword,
			expectedErr:  os.ErrNotExist,
		},
		"non-existent password file": {
			usernameFile: testUsername,
			passwordFile: "nonexistent-password.txt",
			expectedErr:  os.ErrNotExist,
		},
	} {
		t.Run(name, func(t *testing.T) {
			// Prepare
			t.Parallel()

			// Test
			actual, err := usernameAndPassword(tc.usernameFile, tc.passwordFile)

			// Validate
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
