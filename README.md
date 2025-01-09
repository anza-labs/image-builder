# image-builder

[![GitHub License](https://img.shields.io/github/license/anza-labs/image-builder)][license]
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](code_of_conduct.md)
[![GitHub issues](https://img.shields.io/github/issues/anza-labs/image-builder)](https://github.com/anza-labs/image-builder/issues)
[![GitHub release](https://img.shields.io/github/release/anza-labs/image-builder)](https://GitHub.com/anza-labs/image-builder/releases/)
[![Go Report Card](https://goreportcard.com/badge/github.com/anza-labs/image-builder)](https://goreportcard.com/report/github.com/anza-labs/image-builder)

The `image-builder` project provides a Kubernetes-native solution for automating the creation of customized LinuxKit-based images for deployment environments. It utilizes CRDs (Custom Resource Definitions) to define image specifications. The controller orchestrates the image-building process by creating ConfigMaps and Kubernetes Jobs, managing resources efficiently while updating the status of the custom resources. Built with flexibility and scalability in mind, the Image Builder integrates seamlessly into Kubernetes workflows, supporting extensibility through templates and customizable build parameters.

## License

`image-builder` is licensed under the [Apache-2.0][license].

<!-- Resources -->

[license]: https://github.com/anza-labs/image-builder/blob/main/LICENSE
