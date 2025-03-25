# Easy crosscompile toolkit
# hadolint ignore=DL3006
FROM --platform=$BUILDPLATFORM docker.io/tonistiigi/xx:1.6.1 AS xx

# Build the linuxkit binary
FROM docker.io/library/golang:1.24 AS linuxkit
ENV CGO_ENABLED=0

# renovate: datasource=go depName=github.com/linuxkit/linuxkit
ARG VERSION=v0.0.0-20241001125717-ad95c6fc2e48
RUN go install github.com/linuxkit/linuxkit/src/cmd/linuxkit@${VERSION}

# Build the manager binary
FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.24 AS builder
ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM
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
COPY api/ api/
COPY pkg/ pkg/
COPY internal/ internal/

# Build
ENV CGO_ENABLED=0
RUN xx-go build -trimpath -a -o builder pkg/builder/linuxkit/main.go && \
    xx-verify builder

# Use distroless as minimal base image to package the builder binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/builder .
COPY --from=linuxkit /go/bin/linuxkit .
USER 65532:65532
VOLUME ["/tmp"]

ENTRYPOINT ["/builder"]
