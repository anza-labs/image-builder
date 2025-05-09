# API Reference

## Packages
- [image-builder.anza-labs.dev/v1alpha1](#image-builderanza-labsdevv1alpha1)
- [image-builder.anza-labs.dev/v1alpha2](#image-builderanza-labsdevv1alpha2)
- [image-builder.anza-labs.dev/v1beta1](#image-builderanza-labsdevv1beta1)


## image-builder.anza-labs.dev/v1alpha1

Package v1alpha1 contains API Schema definitions for the  v1alpha1 API group

Deprecated: Due to breaking changes, v1alpha2 is a new default, and this version will be removed in upcoming releases.


### Resource Types
- [Image](#image)
- [ImageList](#imagelist)



#### Image



Image is the Schema for the images API.


Deprecated: Due to breaking changes, v1alpha2 is a new default, and this version will be removed in upcoming releases.



_Appears in:_
- [ImageList](#imagelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `image-builder.anza-labs.dev/v1alpha1` | | |
| `kind` _string_ | `Image` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[ImageSpec](#imagespec)_ |  |  |  |
| `status` _[ImageStatus](#imagestatus)_ |  |  |  |


#### ImageList



ImageList contains a list of Image.


Deprecated: Due to breaking changes, v1alpha2 is a new default, and this version will be removed in upcoming releases.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `image-builder.anza-labs.dev/v1alpha1` | | |
| `kind` _string_ | `ImageList` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[Image](#image) array_ |  |  |  |


#### ImageSpec



ImageSpec defines the desired state of Image.



_Appears in:_
- [Image](#image)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `builderImage` _string_ | BuilderImage indicates the container image to use for the Builder job. |  |  |
| `builderVerbosity` _integer_ | BuilderVerbosity specifies log verbosity of the builder. | 4 | Maximum: 10 <br />Minimum: 0 <br /> |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)_ | Resources describe the compute resource requirements. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#affinity-v1-core)_ | Affinity specifies the scheduling constraints for Pods. |  |  |
| `format` _string_ | Format specifies the image format. |  | Enum: [aws docker dynamic-vhd gcp iso-bios iso-efi iso-efi-initrd kernel+initrd kernel+iso kernel+squashfs qcow2-bios qcow2-efi raw-bios raw-efi rpi3 tar tar-kernel-initrd vhd vmdk] <br /> |
| `configuration` _string_ | Configuration is a YAML formatted Linuxkit config. |  |  |
| `result` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | Result is a local reference that lists downloadable objects, that are results of the image building.<br />Defaults to the Image.Metadata.Name. |  |  |
| `bucketCredentials` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | BucketCredentials is a reference to the credentials for S3, where the image will be stored. |  |  |


#### ImageStatus



ImageStatus defines the observed state of Image.



_Appears in:_
- [Image](#image)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates whether the image is ready. |  |  |



## image-builder.anza-labs.dev/v1alpha2

Package v1alpha2 contains API Schema definitions for the image-builder v1alpha2 API group

Deprecated: Due to breaking changes, v1beta1 is a new default, and this version will be removed in upcoming releases.


### Resource Types
- [Image](#image)
- [ImageList](#imagelist)



#### AdditionalData



AdditionalData represents additional data sources for image building.



_Appears in:_
- [ImageSpec](#imagespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name specifies unique name for the additional data. |  |  |
| `volumeMountPoint` _string_ | VolumeMountPoint specifies the path where this data should be mounted. |  |  |
| `configMap` _[ConfigMapVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#configmapvolumesource-v1-core)_ | ConfigMap specifies a ConfigMap as a data source. |  |  |
| `secret` _[SecretVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#secretvolumesource-v1-core)_ | Secret specifies a Secret as a data source. |  |  |
| `image` _[ImageVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#imagevolumesource-v1-core)_ | Image specifies a container image as a data source. |  |  |
| `volume` _[PersistentVolumeClaimVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#persistentvolumeclaimvolumesource-v1-core)_ | Volume specifies a PersistentVolumeClaim as a data source. |  |  |
| `bucket` _[BucketDataSource](#bucketdatasource)_ | Bucket specifies an S3 bucket as a data source. |  |  |
| `gitRepository` _[GitRepository](#gitrepository)_ | GitRepository specifies a Git repository as a data source. |  |  |


#### BucketDataSource



BucketDataSource represents an S3 bucket data source.



_Appears in:_
- [AdditionalData](#additionaldata)
- [DataSource](#datasource)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `credentials` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | Credentials is a reference to the credentials for accessing the bucket. |  |  |
| `items` _[KeyToPath](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#keytopath-v1-core) array_ | Items specifies specific items within the bucket to include. |  |  |
| `itemsConfigMap` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | ItemsSecret specifies a Scret mapping item names to object storage keys.<br />Each value should either be a key of the object or follow the format "key = <Presigned URL>",<br />e.g.:<br />	item-1: "path/to/item-1 = <Presigned URL>"<br />	item-2: "path/to/item-2" |  |  |


#### Container







_Appears in:_
- [ImageSpec](#imagespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | Image indicates the container image to use for the init container. |  |  |
| `verbosity` _integer_ | Verbosity specifies the log verbosity level for the container. | 4 | Maximum: 10 <br />Minimum: 0 <br /> |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)_ | Resources describe the compute resource requirements for the builder job. |  |  |


#### DataSource



DataSource defines the available sources for additional data.
Each data source is either used directly as a Volume for the image, or
will be fetched into empty dir shared between init container and the builder.



_Appears in:_
- [AdditionalData](#additionaldata)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `configMap` _[ConfigMapVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#configmapvolumesource-v1-core)_ | ConfigMap specifies a ConfigMap as a data source. |  |  |
| `secret` _[SecretVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#secretvolumesource-v1-core)_ | Secret specifies a Secret as a data source. |  |  |
| `image` _[ImageVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#imagevolumesource-v1-core)_ | Image specifies a container image as a data source. |  |  |
| `volume` _[PersistentVolumeClaimVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#persistentvolumeclaimvolumesource-v1-core)_ | Volume specifies a PersistentVolumeClaim as a data source. |  |  |
| `bucket` _[BucketDataSource](#bucketdatasource)_ | Bucket specifies an S3 bucket as a data source. |  |  |
| `gitRepository` _[GitRepository](#gitrepository)_ | GitRepository specifies a Git repository as a data source. |  |  |


#### GitRepository



GitRepository represents a Git repository data source.



_Appears in:_
- [AdditionalData](#additionaldata)
- [DataSource](#datasource)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `repository` _string_ | Repository specifies the URL of the Git repository. |  |  |
| `ref` _string_ | Ref specifies the branch, tag, or commit hash to be used from the Git repository. |  |  |
| `credentials` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | Credentials specifies the credentials for accessing the repository.<br />Secret must be one of the following types:<br />	- "kubernetes.io/basic-auth" with "username" and "password" fields;<br />	- "kubernetes.io/ssh-auth" with "ssh-privatekey" field;<br />	- "Opaque" with "gitconfig" field. |  |  |


#### Image



Image represents the schema for the images API.



_Appears in:_
- [ImageList](#imagelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `image-builder.anza-labs.dev/v1alpha2` | | |
| `kind` _string_ | `Image` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[ImageSpec](#imagespec)_ |  |  |  |
| `status` _[ImageStatus](#imagestatus)_ |  |  |  |


#### ImageList



ImageList contains a list of Image resources.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `image-builder.anza-labs.dev/v1alpha2` | | |
| `kind` _string_ | `ImageList` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[Image](#image) array_ |  |  |  |


#### ImageSpec



ImageSpec defines the desired state of an Image resource.



_Appears in:_
- [Image](#image)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `builder` _[Container](#container)_ | Builder specifies the parameters for the main container configuration. |  |  |
| `objFetcher` _[Container](#container)_ | ObjFetcher specifies the parameters for the Object Fetcher init container configuration. |  |  |
| `gitFetcher` _[Container](#container)_ | GitFetcher specifies the parameters for the Git Fetcher init container configuration. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#affinity-v1-core)_ | Affinity specifies the scheduling constraints for Pods running the builder job. |  |  |
| `format` _string_ | Format specifies the output image format. |  | Enum: [aws docker dynamic-vhd gcp iso-bios iso-efi iso-efi-initrd kernel+initrd kernel+iso kernel+squashfs qcow2-bios qcow2-efi raw-bios raw-efi rpi3 tar tar-kernel-initrd vhd vmdk] <br /> |
| `configuration` _string_ | Configuration is a YAML-formatted Linuxkit configuration. |  |  |
| `result` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | Result is a reference to the local object containing downloadable build results.<br />Defaults to the Image.Metadata.Name if not specified. |  |  |
| `bucketCredentials` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | BucketCredentials is a reference to the credentials used for storing the image in S3. |  |  |
| `additionalData` _[AdditionalData](#additionaldata) array_ | AdditionalData specifies additional data sources required for building the image. |  |  |


#### ImageStatus



ImageStatus defines the observed state of an Image resource.



_Appears in:_
- [Image](#image)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates whether the image has been successfully built. |  |  |



## image-builder.anza-labs.dev/v1beta1

Package v1beta1 contains API Schema definitions for the image-builder v1beta1 API group.

### Resource Types
- [LinuxKit](#linuxkit)
- [LinuxKitList](#linuxkitlist)
- [Mkosi](#mkosi)
- [MkosiList](#mkosilist)



#### AdditionalData



AdditionalData represents additional data sources for image building.



_Appears in:_
- [LinuxKitSpec](#linuxkitspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name specifies unique name for the additional data. |  |  |
| `volumeMountPoint` _string_ | VolumeMountPoint specifies the path where this data should be mounted. |  |  |
| `configMap` _[ConfigMapVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#configmapvolumesource-v1-core)_ | ConfigMap specifies a ConfigMap as a data source. |  |  |
| `secret` _[SecretVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#secretvolumesource-v1-core)_ | Secret specifies a Secret as a data source. |  |  |
| `image` _[ImageVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#imagevolumesource-v1-core)_ | Image specifies a container image as a data source. |  |  |
| `volume` _[PersistentVolumeClaimVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#persistentvolumeclaimvolumesource-v1-core)_ | Volume specifies a PersistentVolumeClaim as a data source. |  |  |
| `bucket` _[BucketDataSource](#bucketdatasource)_ | Bucket specifies an S3 bucket as a data source. |  |  |
| `gitRepository` _[GitRepository](#gitrepository)_ | GitRepository specifies a Git repository as a data source. |  |  |


#### BucketDataSource



BucketDataSource represents an S3 bucket data source.



_Appears in:_
- [AdditionalData](#additionaldata)
- [DataSource](#datasource)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `credentials` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | Credentials is a reference to the credentials for accessing the bucket. |  |  |
| `items` _[KeyToPath](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#keytopath-v1-core) array_ | Items specifies specific items within the bucket to include. |  |  |
| `itemsConfigMap` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | ItemsSecret specifies a Scret mapping item names to object storage keys.<br />Each value should either be a key of the object or follow the format "key = <Presigned URL>",<br />e.g.:<br />	item-1: "path/to/item-1 = <Presigned URL>"<br />	item-2: "path/to/item-2" |  |  |


#### Container







_Appears in:_
- [LinuxKitSpec](#linuxkitspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | Image indicates the container image to use for the init container. |  |  |
| `verbosity` _integer_ | Verbosity specifies the log verbosity level for the container. | 4 | Maximum: 10 <br />Minimum: 0 <br /> |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)_ | Resources describe the compute resource requirements for the builder job. |  |  |


#### DataSource



DataSource defines the available sources for additional data.
Each data source is either used directly as a Volume for the image, or
will be fetched into empty dir shared between init container and the builder.



_Appears in:_
- [AdditionalData](#additionaldata)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `configMap` _[ConfigMapVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#configmapvolumesource-v1-core)_ | ConfigMap specifies a ConfigMap as a data source. |  |  |
| `secret` _[SecretVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#secretvolumesource-v1-core)_ | Secret specifies a Secret as a data source. |  |  |
| `image` _[ImageVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#imagevolumesource-v1-core)_ | Image specifies a container image as a data source. |  |  |
| `volume` _[PersistentVolumeClaimVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#persistentvolumeclaimvolumesource-v1-core)_ | Volume specifies a PersistentVolumeClaim as a data source. |  |  |
| `bucket` _[BucketDataSource](#bucketdatasource)_ | Bucket specifies an S3 bucket as a data source. |  |  |
| `gitRepository` _[GitRepository](#gitrepository)_ | GitRepository specifies a Git repository as a data source. |  |  |


#### GitRepository



GitRepository represents a Git repository data source.



_Appears in:_
- [AdditionalData](#additionaldata)
- [DataSource](#datasource)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `repository` _string_ | Repository specifies the URL of the Git repository. |  |  |
| `ref` _string_ | Ref specifies the branch, tag, or commit hash to be used from the Git repository. |  |  |
| `credentials` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | Credentials specifies the credentials for accessing the repository.<br />Secret must be one of the following types:<br />	- "kubernetes.io/basic-auth" with "username" and "password" fields;<br />	- "kubernetes.io/ssh-auth" with "ssh-privatekey" field;<br />	- "Opaque" with "gitconfig" field. |  |  |


#### LinuxKit



LinuxKit is the Schema for the linuxkits API.



_Appears in:_
- [LinuxKitList](#linuxkitlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `image-builder.anza-labs.dev/v1beta1` | | |
| `kind` _string_ | `LinuxKit` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[LinuxKitSpec](#linuxkitspec)_ |  |  |  |
| `status` _[LinuxKitStatus](#linuxkitstatus)_ |  |  |  |


#### LinuxKitList



LinuxKitList contains a list of LinuxKit.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `image-builder.anza-labs.dev/v1beta1` | | |
| `kind` _string_ | `LinuxKitList` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[LinuxKit](#linuxkit) array_ |  |  |  |


#### LinuxKitSpec



LinuxKitSpec defines the desired state of an LinuxKit resource.



_Appears in:_
- [LinuxKit](#linuxkit)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `builder` _[Container](#container)_ | Builder specifies the parameters for the main container configuration. |  |  |
| `objFetcher` _[Container](#container)_ | ObjFetcher specifies the parameters for the Object Fetcher init container configuration. |  |  |
| `gitFetcher` _[Container](#container)_ | GitFetcher specifies the parameters for the Git Fetcher init container configuration. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#affinity-v1-core)_ | Affinity specifies the scheduling constraints for Pods running the builder job. |  |  |
| `format` _string_ | Format specifies the output image format. |  | Enum: [aws docker dynamic-vhd gcp iso-bios iso-efi iso-efi-initrd kernel+initrd kernel+iso kernel+squashfs qcow2-bios qcow2-efi raw-bios raw-efi rpi3 tar tar-kernel-initrd vhd vmdk] <br /> |
| `configuration` _string_ | Configuration is a YAML-formatted Linuxkit configuration. |  |  |
| `result` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | Result is a reference to the local object containing downloadable build results.<br />Defaults to the Image.Metadata.Name if not specified. |  |  |
| `bucketCredentials` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | BucketCredentials is a reference to the credentials used for storing the image in S3. |  |  |
| `additionalData` _[AdditionalData](#additionaldata) array_ | AdditionalData specifies additional data sources required for building the image. |  |  |


#### LinuxKitStatus



LinuxKitStatus defines the observed state of an Image resource.



_Appears in:_
- [LinuxKit](#linuxkit)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates whether the image has been successfully built. |  |  |


#### Mkosi



Mkosi is the Schema for the mkosis API.



_Appears in:_
- [MkosiList](#mkosilist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `image-builder.anza-labs.dev/v1beta1` | | |
| `kind` _string_ | `Mkosi` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[MkosiSpec](#mkosispec)_ |  |  |  |
| `status` _[MkosiStatus](#mkosistatus)_ |  |  |  |


#### MkosiList



MkosiList contains a list of Mkosi.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `image-builder.anza-labs.dev/v1beta1` | | |
| `kind` _string_ | `MkosiList` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[Mkosi](#mkosi) array_ |  |  |  |


#### MkosiSpec



MkosiSpec defines the desired state of Mkosi.



_Appears in:_
- [Mkosi](#mkosi)



#### MkosiStatus



MkosiStatus defines the observed state of Mkosi.



_Appears in:_
- [Mkosi](#mkosi)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates whether the image has been successfully built. |  |  |


