# syntax=docker/dockerfile:experimental
ARG GO_VERSION=1.22
ARG DEBIAN_VERSION=11-slim
FROM golang:${GO_VERSION} as base

FROM golang:${GO_VERSION} as releaser

RUN --mount=type=cache,target=/go/pkg/cache \
    go install github.com/goreleaser/goreleaser@latest
