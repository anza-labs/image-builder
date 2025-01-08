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
	"encoding/json"
	"fmt"
	"path/filepath"

	imagebuilderv1alpha2 "github.com/anza-labs/image-builder/api/v1alpha2"
	"github.com/anza-labs/image-builder/internal/fetcherconfig"
	"github.com/anza-labs/image-builder/internal/naming"
	"github.com/anza-labs/image-builder/version"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func Role(image *imagebuilderv1alpha2.Image) *rbacv1.Role {
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

func RoleBinding(image *imagebuilderv1alpha2.Image) *rbacv1.RoleBinding {
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

func ServiceAccount(image *imagebuilderv1alpha2.Image) *corev1.ServiceAccount {
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

func config(image *imagebuilderv1alpha2.Image) (string, error) {
	data := fetcherconfig.Config{}
	for _, ad := range image.Spec.AdditionalData {
		if f := NewConfigMapEntryFrom(ad); f != nil {
			data.Fetchers = append(data.Fetchers, *f)
		}
	}

	b, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("unable to encode configuration: %w", err)
	}

	return string(b), nil
}

func InitConfigMap(image *imagebuilderv1alpha2.Image) (*corev1.ConfigMap, error) {
	config, err := config(image)
	if err != nil {
		return nil, err
	}
	h := fmt.Sprintf("%x", sha256.Sum256([]byte(config)))

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      naming.ConfigMap(image.Name, h), // TODO
			Namespace: image.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       image.Name,
				"app.kubernetes.io/managed-by": "image-builder",
			},
		},
		Data: map[string]string{
			"fetcher.json": config,
		},
	}, nil
}

func ConfigMap(image *imagebuilderv1alpha2.Image) *corev1.ConfigMap {
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

func Job(image *imagebuilderv1alpha2.Image) (*batchv1.Job, error) {
	affinity := image.Spec.Affinity

	outputSecret := image.Spec.Result
	if outputSecret.Name == "" {
		outputSecret.Name = image.Name
	}

	volumes := DefaultVolumes(image)
	volumeMounts := []corev1.VolumeMount{}

	initVolumeMounts := []corev1.VolumeMount{}

	for _, d := range image.Spec.AdditionalData {
		vo := NewVolumeFrom(d)
		volumes = append(volumes, vo.volumes...)
		volumeMounts = append(volumeMounts, vo.volumeMount)
		initVolumeMounts = append(initVolumeMounts, vo.initVolumeMounts...)
	}

	config, err := config(image)
	if err != nil {
		return nil, err
	}
	h := fmt.Sprintf("%x", sha256.Sum256([]byte(config)))
	fetcherCM := naming.ConfigMap(image.Name, h)
	volumes = append(volumes, corev1.Volume{
		Name: fetcherCM,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: fetcherCM,
				},
			},
		},
	})
	initVolumeMounts = append(initVolumeMounts, corev1.VolumeMount{
		Name:      fetcherCM,
		ReadOnly:  true,
		MountPath: "/etc/fetcher/config.json",
	})

	containers := []corev1.Container{
		Container(image, volumeMounts...),
	}

	initContainers := []corev1.Container{
		InitCointainer(image.Spec.GitFetcher, "gitfetcher", initVolumeMounts...),
		InitCointainer(image.Spec.ObjFetcher, "objfetcher", initVolumeMounts...),
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
					InitContainers:     initContainers,
					Containers:         containers,
					Volumes:            volumes,
					Affinity:           affinity,
					ServiceAccountName: image.Name,
					RestartPolicy:      corev1.RestartPolicyNever,
				},
			},
		},
	}, nil
}

func DefaultVolumes(image *imagebuilderv1alpha2.Image) []corev1.Volume {
	bucketCredentials := image.Spec.BucketCredentials
	config := image.Spec.Configuration
	h := fmt.Sprintf("%x", sha256.Sum256([]byte(config)))

	return []corev1.Volume{
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
	}
}

func Container(image *imagebuilderv1alpha2.Image, extraVolumeMounts ...corev1.VolumeMount) corev1.Container {
	outputSecret := image.Spec.Result
	if outputSecret.Name == "" {
		outputSecret.Name = image.Name
	}

	containerImage := image.Spec.Builder.Image
	if containerImage == "" {
		containerImage = fmt.Sprintf("%s/image-builder:%s", version.OCIRepository, version.Version)
	}

	format := image.Spec.Format
	resources := image.Spec.Builder.Resources
	verbosity := image.Spec.Builder.Verbosity

	volumeMounts := []corev1.VolumeMount{
		{Name: "bucket-credentials", MountPath: "/credentials"},
		{Name: "config", MountPath: "/config"},
		{Name: "temp", MountPath: "/tmp"},
	}
	volumeMounts = append(volumeMounts, extraVolumeMounts...)

	return corev1.Container{
		Name:  "builder",
		Image: containerImage,
		Args: []string{
			fmt.Sprintf("--v=%d", verbosity),
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
		VolumeMounts: volumeMounts,
		Resources:    resources,
	}
}

func InitCointainer(
	ctr imagebuilderv1alpha2.Container,
	name string,
	extraVolumeMounts ...corev1.VolumeMount,
) corev1.Container {
	containerImage := ctr.Image
	if containerImage == "" {
		containerImage = fmt.Sprintf("%s/image-builder-init-%s:%s", version.OCIRepository, name, version.Version)
	}

	resources := ctr.Resources
	verbosity := ctr.Verbosity

	volumeMounts := []corev1.VolumeMount{}
	volumeMounts = append(volumeMounts, extraVolumeMounts...)

	return corev1.Container{
		Name:  naming.InitCointainer(name),
		Image: containerImage,
		Args: []string{
			fmt.Sprintf("--v=%d", verbosity),
		},
		VolumeMounts: volumeMounts,
		Resources:    resources,
		Env: []corev1.EnvVar{
			{Name: "K8S_NAME", ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.name"},
			}},
			{Name: "K8S_NAMESPACE", ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.namespace"},
			}},
			{Name: "FETCHER_CONFIG", Value: "/config/fetcher.json"},
		},
	}
}

type volumeOpts struct {
	volumes          []corev1.Volume
	volumeMount      corev1.VolumeMount
	initVolumeMounts []corev1.VolumeMount
}

func NewVolumeFrom(data imagebuilderv1alpha2.AdditionalData) volumeOpts {
	var source corev1.VolumeSource

	vo := volumeOpts{
		volumeMount: corev1.VolumeMount{
			Name:      data.Name,
			MountPath: data.VolumeMountPoint,
		},
	}

	if data.DataSource.Bucket != nil {
		source = corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{
				Medium: "",
			},
		}

		if data.DataSource.Bucket.ItemsSecret != nil {
			items := naming.Volume("%s-%s", data.Name, "items")

			vo.volumes = append(vo.volumes, corev1.Volume{
				Name: items,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: data.Bucket.ItemsSecret.Name,
					},
				},
			})
			vo.initVolumeMounts = append(vo.initVolumeMounts, corev1.VolumeMount{
				Name:      items,
				MountPath: filepath.Join("/etc/objfetcher", items),
			})
		}

		objCreds := naming.Volume("%s-%s", data.Name, "objcreds")
		vo.volumes = append(vo.volumes, corev1.Volume{
			Name: objCreds,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: data.Bucket.Credentials.Name,
				},
			},
		})
		vo.initVolumeMounts = append(vo.initVolumeMounts, corev1.VolumeMount{
			Name:      objCreds,
			MountPath: filepath.Join("/etc/objfetcher", objCreds),
		})
	}

	if data.DataSource.ConfigMap != nil {
		source = corev1.VolumeSource{
			ConfigMap: data.DataSource.ConfigMap,
		}
	}

	if data.DataSource.GitRepository != nil {
		source = corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{
				Medium: "",
			},
		}

		if data.GitRepository.Credentials != nil {
			gitCreds := naming.Volume("%s-%s", data.Name, "gitcreds")
			vo.volumes = append(vo.volumes,
				corev1.Volume{
					Name: gitCreds,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: data.GitRepository.Credentials.Name,
						},
					},
				},
			)
			vo.initVolumeMounts = append(vo.initVolumeMounts,
				corev1.VolumeMount{
					Name:      gitCreds,
					MountPath: filepath.Join("/etc/gitfetcher", gitCreds),
				},
			)
		}
	}

	if data.DataSource.Image != nil {
		source = corev1.VolumeSource{
			Image: data.DataSource.Image,
		}
	}

	if data.DataSource.Secret != nil {
		source = corev1.VolumeSource{
			Secret: data.DataSource.Secret,
		}
	}

	if data.DataSource.Volume != nil {
		source = corev1.VolumeSource{
			PersistentVolumeClaim: data.DataSource.Volume,
		}
	}

	vo.volumes = append(vo.volumes, corev1.Volume{
		Name:         data.Name,
		VolumeSource: source,
	})
	vo.initVolumeMounts = append(vo.initVolumeMounts, vo.volumeMount)

	return vo
}

func mode(m *int32) int32 {
	if m != nil {
		return *m
	}
	return 0o755
}

func NewConfigMapEntryFrom(data imagebuilderv1alpha2.AdditionalData) *fetcherconfig.Fetcher {
	config := &fetcherconfig.Fetcher{}

	if data.DataSource.Bucket != nil {
		objCreds := naming.Volume("%s-%s", data.Name, "objcreds")
		config.ObjFetcher = &fetcherconfig.ObjFetcher{
			MountPoint:      data.VolumeMountPoint,
			CredentialsPath: filepath.Join("/etc/objfetcher", objCreds),
		}

		if data.DataSource.Bucket.ItemsSecret != nil {
			items := naming.Volume("%s-%s", data.Name, "items")
			config.ObjFetcher.KeysPath = filepath.Join("/etc/objfetcher", items)
		}

		for _, kv := range data.DataSource.Bucket.Items {
			if config.ObjFetcher.Keys == nil {
				config.ObjFetcher.Keys = make(map[string]fetcherconfig.File)
			}
			config.ObjFetcher.Keys[kv.Key] = fetcherconfig.File{
				Path: kv.Path,
				Mode: mode(kv.Mode),
			}
		}

		return config
	}

	if data.DataSource.GitRepository != nil {
		gitCreds := naming.Volume("%s-%s", data.Name, "gitcreds")
		config.GitFetcher = &fetcherconfig.GitFetcher{
			MountPoint:      data.VolumeMountPoint,
			CredentialsPath: filepath.Join("/etc/gitfetcher", gitCreds),
			Repository:      data.DataSource.GitRepository.Repository,
			Ref:             data.DataSource.GitRepository.Ref,
		}

		return config
	}

	// otherwise, no-op
	return nil
}
