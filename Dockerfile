# syntax=docker/dockerfile:1.2

FROM --platform=$BUILDPLATFORM golang:1-alpine AS builder
WORKDIR /w
ENV CGO_ENABLED=0
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
RUN \
  --mount=type=cache,target=/var/cache/apk ln -vs /var/cache/apk /etc/apk/cache && \
    set -ux \
 && apk update \
 && apk add git \
 && git init
COPY go.??? .
RUN \
  --mount=type=cache,target=/go/pkg/mod \
    set -ux \
 && git add -A . \
 && go mod download \
 && git --no-pager diff && [[ 0 -eq $(git --no-pager diff --name-only | wc -l) ]]
COPY . .

RUN \
  --mount=type=cache,target=/go/pkg/mod \
    set -ux \
 && git add -A . \
 && go mod tidy \
 && go mod verify \
 && git --no-pager diff && [[ 0 -eq $(git --no-pager diff --name-only | wc -l) ]]

RUN \
  --mount=type=cache,target=/go/pkg/mod \
    set -ux \
 && git add -A . \
 && go fmt ./... \
 && git --no-pager diff && [[ 0 -eq $(git --no-pager diff --name-only | wc -l) ]]

RUN \
  --mount=type=cache,target=/go/pkg/mod \
    set -ux \
 && git add -A . \
 && go vet ./... \
 && git --no-pager diff && [[ 0 -eq $(git --no-pager diff --name-only | wc -l) ]]

RUN \
  --mount=type=cache,target=/go/pkg/mod \
    set -ux \
 && git add -A . \
 && go test ./... \
 && git --no-pager diff && [[ 0 -eq $(git --no-pager diff --name-only | wc -l) ]]

RUN \
  --mount=type=cache,target=/go/pkg/mod \
    set -ux \
 && go build -o lck -ldflags '-s -w' cmd/locked.go

FROM --platform=$BUILDPLATFORM scratch AS binaries
COPY --from=builder /w/lck* /
