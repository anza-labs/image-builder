---
title: Quick Start
weight: 1
---

## Minimal image

The following example demonstrates how to create a minimal image specification using the Image Builder API. This setup leverages the [Container Object Storage Interface (COSI)][cosi] for managing bucket credentials.


!!! note
    You can use any of the object storage providers, that support S3 protocol and has [COSI (v1alpha1)][cosi] driver.

    Example of such drivers are:

      * [Ceph COSI](https://github.com/ceph/ceph-cosi)
      * [Linode COSI Driver](https://github.com/linode/linode-cosi-driver)

### Example Configuration

```yaml
apiVersion: image-builder.anza-labs.dev/v1alpha2
kind: Image
metadata:
  name: minimal
spec:
  format: 'kernel+initrd'
  configuration: |
    kernel:
      image: linuxkit/kernel:6.6.13
      cmdline: "console=tty0 console=ttyS0 console=ttyAMA0"
    init:
      - linuxkit/init:e120ea2a30d906bd1ee1874973d6e4b1403b5ca3
      - linuxkit/runc:6062483d748609d505f2bcde4e52ee64a3329f5f
      - linuxkit/containerd:39301e7312f13eedf19bd5d5551af7b37001d435
    onboot:
      - name: dhcpcd
        image: linuxkit/dhcpcd:e9e3580f2de00e73e7b316a007186d22fea056ee
        command: ["/sbin/dhcpcd", "--nobackground", "-f", "/dhcpcd.conf", "-1"]
    services:
      - name: getty
        image: linuxkit/getty:5d86a2ce2d890c14ab66b13638dcadf74f29218b
        env:
        - INSECURE=true
  bucketCredentials:
    name: s3-credentials
```

### Applying the Configuration

Save the above configuration into a file named `minimal-image.yaml`. You can then apply it using `kubectl`:

```sh
kubectl apply -f minimal-image.yaml
```

Once applied, you can verify the status and access the generated secrets as follows:

```
$ kubectl get images.image-builder.anza-labs.dev minimal
NAME      READY
minimal   true
```

To view the created secrets:

```
$ kubectl get secrets minimal
NAME      TYPE     DATA   AGE
minimal   Opaque   3      3m36s
```

To retrieve the key and presigned URL from the secret:

```
$ kubectl get secrets minimal -o=json | jq -r '.data."<Output Name>"' | base64 -d
<Object Key> = <Presigned URL>
```

You can then use e.g. [mc](https://min.io/docs/minio/linux/reference/minio-mc.html) to fetch the objects:

```
$ mc ls <alias>/<bucket>/default/minimal/kernel-initrd/
[2024-12-30 18:57:13 CET]    42B STANDARD image-cmdline
[2024-12-30 19:00:48 CET]  77MiB STANDARD image-initrd-img
[2024-12-30 19:01:12 CET] 8.7MiB STANDARD image-kernel
```

[cosi]: https://github.com/kubernetes-sigs/container-object-storage-interface
