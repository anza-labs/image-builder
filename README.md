# image-builder

[![GitHub License](https://img.shields.io/github/license/anza-labs/image-builder)][license]
[![GitHub issues](https://img.shields.io/github/issues/anza-labs/image-builder)](https://github.com/anza-labs/image-builder/issues)
[![GitHub release](https://img.shields.io/github/release/anza-labs/image-builder)](https://GitHub.com/anza-labs/image-builder/releases/)

The `image-builder` project provides a Kubernetes-native solution for automating the creation of customized LinuxKit-based images for deployment environments. It utilizes CRDs (Custom Resource Definitions) to define image specifications. The controller orchestrates the image-building process by creating ConfigMaps and Kubernetes Jobs, managing resources efficiently while updating the status of the custom resources. Built with flexibility and scalability in mind, the Image Builder integrates seamlessly into Kubernetes workflows, supporting extensibility through templates and customizable build parameters.

## License

`image-builder` is licensed under the [Apache-2.0][license].

<!-- Resources -->

[license]: https://github.com/anza-labs/image-builder/blob/main/LICENSE
