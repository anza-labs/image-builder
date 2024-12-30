# API Reference

## Packages
- [image-builder.anza-labs.dev/v1alpha1](#image-builderanza-labsdevv1alpha1)


## image-builder.anza-labs.dev/v1alpha1

Package v1alpha1 contains API Schema definitions for the  v1alpha1 API group

### Resource Types
- [Image](#image)
- [ImageList](#imagelist)



#### Image



Image is the Schema for the images API.



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


