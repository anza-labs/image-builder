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

package template

import (
	"bytes"
	"context"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestKubeFuncs(t *testing.T) {
	scheme := runtime.NewScheme()
	assert.NoError(t, corev1.AddToScheme(scheme))

	testNamespace := "test-namespace"

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: testNamespace,
		},
		StringData: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: testNamespace,
		},
		Data: map[string]string{
			"keyA": "valueA",
			"keyB": "valueB",
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret, configMap).Build()

	for name, tt := range map[string]struct {
		function       interface{}
		args           []interface{}
		expectedResult interface{}
	}{
		"fetchSecretData - valid": {
			function: fetchSecretData,
			args:     []interface{}{"test-secret"},
			expectedResult: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		"fetchSecretKey - valid key": {
			function:       fetchSecretKey,
			args:           []interface{}{"test-secret", "key1"},
			expectedResult: "value1",
		},
		"fetchConfigMapData - valid": {
			function: fetchConfigMapData,
			args:     []interface{}{"test-configmap"},
			expectedResult: map[string]string{
				"keyA": "valueA",
				"keyB": "valueB",
			},
		},
		"fetchConfigMapKey - valid key": {
			function:       fetchConfigMapKey,
			args:           []interface{}{"test-configmap", "keyA"},
			expectedResult: "valueA",
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			var result interface{}

			switch closure := tt.function.(type) {
			case func(ctx context.Context, namespace string, cli client.Client) dataFunc:
				fn := closure(ctx, testNamespace, fakeClient)

				assert.NotPanics(t, func() {
					result = fn(tt.args[0].(string))
					assert.Equal(t, tt.expectedResult, result)
				})

			case func(ctx context.Context, namespace string, cli client.Client) dataKeyFunc:
				fn := closure(ctx, testNamespace, fakeClient)
				assert.NotPanics(t, func() {
					result = fn(tt.args[0].(string), tt.args[1].(string))
					assert.Equal(t, tt.expectedResult, result)
				})
			}
		})
	}
}

func TestTemplateFuncs(t *testing.T) {
	scheme := runtime.NewScheme()
	assert.NoError(t, corev1.AddToScheme(scheme))

	testNamespace := "test-namespace"

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: testNamespace,
		},
		StringData: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: testNamespace,
		},
		Data: map[string]string{
			"keyA": "valueA",
			"keyB": "valueB",
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret, configMap).Build()

	ctx := context.Background()
	funcs := Funcs(ctx, testNamespace, fakeClient)

	for name, tc := range map[string]struct {
		template string
		expected string
		hasError bool
	}{
		"Fetch secret key": {
			template: `{{ fetchSecretKey "test-secret" "key1" }}`,
			expected: "value1",
		},
		"Fetch configmap key": {
			template: `{{ fetchConfigMapKey "test-configmap" "keyA" }}`,
			expected: "valueA",
		},
		"Fetch non-existent secret": {
			template: `{{ fetchSecretData "foo-secret" }}`,
			hasError: true,
		},
		"Fetch all secret data": {
			template: `{{ range $k, $v := fetchSecretData "test-secret" }}{{$k}}={{$v}},{{ end }}`,
			expected: "key1=value1,key2=value2,",
		},
		"Fetch all configmap data": {
			template: `{{ range $k, $v := fetchConfigMapData "test-configmap" }}{{$k}}={{$v}},{{ end }}`,
			expected: "keyA=valueA,keyB=valueB,",
		},
	} {
		t.Run(name, func(t *testing.T) {
			tmpl, err := template.New("test").Funcs(funcs).Parse(tc.template)
			assert.NoError(t, err)

			var buf bytes.Buffer
			err = tmpl.Execute(&buf, nil)

			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, buf.String())
			}
		})
	}
}
