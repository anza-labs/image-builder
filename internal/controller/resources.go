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

	anzalabsdevv1alpha1 "github.com/anza-labs/image-builder/api/v1alpha1"
	"github.com/anza-labs/image-builder/internal/naming"
	"github.com/anza-labs/image-builder/version"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func Role(image *anzalabsdevv1alpha1.Image) *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      image.Name,
			Namespace: image.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       image.Name,
				"app.kubernetes.io/managed-by": "image-builder",
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"create", "delete", "get", "list", "patch", "update", "watch"},
			},
		},
	}
}

func RoleBinding(image *anzalabsdevv1alpha1.Image) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      image.Name,
			Namespace: image.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       image.Name,
				"app.kubernetes.io/managed-by": "image-builder",
			},
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      image.Name,
				Namespace: image.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     image.Name,
		},
	}
}

func ServiceAccount(image *anzalabsdevv1alpha1.Image) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      image.Name,
			Namespace: image.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       image.Name,
				"app.kubernetes.io/managed-by": "image-builder",
			},
		},
		AutomountServiceAccountToken: ptr.To(true),
	}
}

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
			"image.yaml": config,
		},
	}
}

func Job(image *anzalabsdevv1alpha1.Image) *batchv1.Job {
	outputSecret := image.Spec.Result
	bucketCredentials := image.Spec.BucketCredentials
	format := image.Spec.Format
	containerImage := image.Spec.BuilderImage
	affinity := image.Spec.Affinity
	resources := image.Spec.Resources

	config := image.Spec.Configuration
	h := fmt.Sprintf("%x", sha256.Sum256([]byte(config)))

	if outputSecret.Name == "" {
		outputSecret.Name = image.Name
	}

	if containerImage == "" {
		containerImage = fmt.Sprintf("%s/image-builder:%s", version.OCIRepository, version.Version)
	}

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      image.Name,
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
							Args: []string{
								"--v=4",
							},
							Env: []corev1.EnvVar{
								{Name: "K8S_NAME", ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.name"},
								}},
								{Name: "K8S_NAMESPACE", ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.namespace"},
								}},
								{Name: "K8S_SECRET_NAME", Value: outputSecret.Name},
								{Name: "LINUXKIT_FORMAT", Value: format},
								{Name: "LINUXKIT_CONFIG", Value: "/config/image.yaml"},
								{Name: "STORAGE_CREDENTIALS", Value: "/credentials/BucketInfo.json"},
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "bucket-credentials", MountPath: "/credentials"},
								{Name: "config", MountPath: "/config"},
								{Name: "temp", MountPath: "/tmp"},
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
						{
							Name: "temp",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{Medium: ""},
							},
						},
					},
					Affinity:           affinity,
					ServiceAccountName: image.Name,
					RestartPolicy:      corev1.RestartPolicyNever,
				},
			},
		},
	}
}
