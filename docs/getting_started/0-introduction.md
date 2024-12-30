---
title: Introduction
weight: 0
---

## Installing the Operator

To install the `image-builder`, run the following commands. This will ensure you're always pulling the latest stable release from the operator’s GitHub repository.

```sh
LATEST="$(curl -s 'https://api.github.com/repos/anza-labs/image-builder/releases/latest' | jq -r '.tag_name')"
kubectl apply -k "https://github.com/anza-labs/image-builder/?ref=${LATEST}"
```

This command:

1. Fetches the latest release tag using the GitHub API.
1. Applies the corresponding version of the `image-builder` to your Kubernetes cluster using `kubectl`.

Once installed, the operator will begin monitoring the appropriate resources in your cluster based on the CRDs defined.
