---
title: Introduction
weight: 0
---

## Installing the Operator

To install the `image-builder`, run the following commands. This will ensure you're always pulling the latest stable release from the operator’s GitHub repository.

!!! note
    This software requires Cert Manager to manage SSL/TLS certificates. Cert Manager automates the creation, renewal, and management of certificates within Kubernetes clusters.

    For installation instructions, please visit the official documentation: [Cert Manager Installation](https://cert-manager.io/docs/installation/).

Let me know if you'd like to modify or expand on any details!
```sh
LATEST="$(curl -s 'https://api.github.com/repos/anza-labs/image-builder/releases/latest' | jq -r '.tag_name')"
kubectl apply -k "https://github.com/anza-labs/image-builder/?ref=${LATEST}"
```

This command:

1. Fetches the latest release tag using the GitHub API.
1. Applies the corresponding version of the `image-builder` to your Kubernetes cluster using `kubectl`.

Once installed, the operator will begin monitoring the appropriate resources in your cluster based on the CRDs defined.
