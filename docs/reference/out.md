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
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[ImageSpec](#imagespec)_ |  |  |  |
| `status` _[ImageStatus](#imagestatus)_ |  |  |  |


#### ImageList



ImageList contains a list of Image.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `image-builder.anza-labs.dev/v1alpha1` | | |
| `kind` _string_ | `ImageList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[Image](#image) array_ |  |  |  |


#### ImageSpec



ImageSpec defines the desired state of Image.



_Appears in:_
- [Image](#image)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `format` _string_ | Format specifies the image format. |  | Enum: [aws docker dynamic-vhd gcp iso-bios iso-efi iso-efi-initrd kernel+initrd kernel+iso kernel+squashfs qcow2-bios qcow2-efi raw-bios raw-efi rpi3 tar tar-kernel-initrd vhd vmdk] <br /> |
| `configuration` _string_ | Configuration is a YAML formatted Linuxkit config. |  |  |
| `bucketCredentials` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#localobjectreference-v1-core)_ | BucketCredentials is a reference to the credentials for S3, where the image will be stored. |  |  |


#### ImageStatus



ImageStatus defines the observed state of Image.



_Appears in:_
- [Image](#image)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates whether the image is ready. |  |  |
| `downloadURL` _string_ | DownloadURL specifies the URL to download the image. |  |  |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#condition-v1-meta) array_ | Conditions lists the conditions of the image resource. |  |  |


