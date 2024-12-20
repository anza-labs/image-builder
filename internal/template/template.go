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
	"context"
	"os"
	"text/template"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type dataFunc func(name string) map[string]string
type dataKeyFunc func(name, key string) string

func fetchSecretData(ctx context.Context, namespace string, cli client.Client) dataFunc {
	return func(name string) map[string]string {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		secret := &corev1.Secret{}
		err := cli.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, secret)
		if err != nil {
			panic(err)
		}

		return secret.StringData
	}
}

func fetchSecretKey(ctx context.Context, namespace string, cli client.Client) dataKeyFunc {
	return func(name, key string) string {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		secret := &corev1.Secret{}
		err := cli.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, secret)
		if err != nil {
			panic(err)
		}

		return secret.StringData[key]
	}
}

func fetchConfigMapData(ctx context.Context, namespace string, cli client.Client) dataFunc {
	return func(name string) map[string]string {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		cm := &corev1.ConfigMap{}
		err := cli.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, cm)
		if err != nil {
			panic(err)
		}

		return cm.Data
	}
}

func fetchConfigMapKey(ctx context.Context, namespace string, cli client.Client) dataKeyFunc {
	return func(name, key string) string {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		cm := &corev1.ConfigMap{}
		err := cli.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, cm)
		if err != nil {
			panic(err)
		}

		return cm.Data[key]
	}
}

func Funcs(ctx context.Context, namespace string, cli client.Client) template.FuncMap {
	return template.FuncMap{
		"getenv":             os.Getenv,
		"fetchSecretData":    fetchSecretData(ctx, namespace, cli),
		"fetchSecretKey":     fetchSecretKey(ctx, namespace, cli),
		"fetchConfigMapData": fetchConfigMapData(ctx, namespace, cli),
		"fetchConfigMapKey":  fetchConfigMapKey(ctx, namespace, cli),
	}
}
