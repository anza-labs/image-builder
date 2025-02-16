# Easy crosscompile toolkit
# hadolint ignore=DL3006
FROM --platform=$BUILDPLATFORM docker.io/tonistiigi/xx:1.6.1 AS xx

# Build the manager binary
FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.24 AS builder
ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM
ARG VERSION="dev"
ARG OCI_REPOSITORY="ghcr.io/anza-labs"
COPY --from=xx / /

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN xx-go mod download

# Copy the go source
COPY hack/ hack/
COPY version/ version/
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY internal/ internal/

# Build
ENV CGO_ENABLED=0
RUN xx-go build -trimpath -a -o manager \
        -ldflags "-X github.com/anza-labs/image-builder/version.Version=$VERSION -X github.com/anza-labs/image-builder/version.OCIRepository=$OCI_REPOSITORY" \
        cmd/main.go && \
    xx-verify manager

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
