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

package controller

import (
	"crypto/sha256"
	"fmt"
	"os"

	anzalabsdevv1alpha1 "github.com/anza-labs/image-builder/api/v1alpha1"
	"github.com/anza-labs/image-builder/internal/naming"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var serviceAccountName = os.Getenv("K8S_SERVICE_ACCOUNT")

func ConfigMap(image *anzalabsdevv1alpha1.Image) *corev1.ConfigMap {
	config := image.Spec.Configuration
	h := fmt.Sprintf("%x", sha256.Sum256([]byte(config)))

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      naming.ConfigMap(image.Name, h),
			Namespace: image.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       image.Name,
				"app.kubernetes.io/managed-by": "image-builder",
			},
		},
		Data: map[string]string{
			"config": config,
		},
	}
}

func Job(image *anzalabsdevv1alpha1.Image) *batchv1.Job {
	outputSecret := image.Spec.Result
	bucketCredentials := image.Spec.BucketCredentials
	format := image.Spec.Format
	containerImage := image.Spec.BuilderTemplate.Image
	affinity := image.Spec.BuilderTemplate.Affinity
	resources := image.Spec.BuilderTemplate.Resources
	serviceAccount := image.Spec.BuilderTemplate.ServiceAccountName

	config := image.Spec.Configuration
	h := fmt.Sprintf("%x", sha256.Sum256([]byte(config)))

	if serviceAccount == "" {
		serviceAccount = serviceAccountName
	}

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-job", image.Name),
			Namespace: image.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       image.Name,
				"app.kubernetes.io/managed-by": "image-builder",
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "builder",
							Image: containerImage,
							Env: []corev1.EnvVar{
								{Name: "OUTPUT", Value: outputSecret.Name},
								{Name: "LINUXKIT_FORMAT", Value: format},
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "bucket-credentials", MountPath: "/credentials"},
								{Name: "config", MountPath: "/config"},
							},
							Resources: resources,
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "bucket-credentials",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: bucketCredentials.Name,
								},
							},
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: naming.ConfigMap(image.Name, h),
									},
								},
							},
						},
					},
					Affinity:           affinity,
					ServiceAccountName: serviceAccount,
					RestartPolicy:      corev1.RestartPolicyNever,
				},
			},
		},
	}
}
