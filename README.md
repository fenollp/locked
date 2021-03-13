# locked
Universal lockfile (not a daemon)

* `locked [--tracked=REGEXP --] ( FILE )+`
* `locked .../Dockerfile`
* looks for single lines matching `^(#|//) +locked--[^:]+: .+$` anywhere in the file
* looks for single lines matching `^(#|//) +locked +[^ :]+:.+$` anywhere in the file
* applies matchers
	* `OCI-image-FROM`: resolves an OCI image URI from a registry to its sha256
* runs concurrently
* writes changes in-place
* does not alter behavior
	* should remain correct and usable without `locked`
	* should resolve to "today"'s values

## Examples

### Dockerfile
Given the Dockerfile
```dockerfile
# locked {
# locked using:OCI-image-FROM
# locked   and:track=golang:1.16.2-alpine3.12
# locked   and:semver=1.16.2~1
FROM golang:1.16.2-alpine3.12 AS builder
# locked }
RUN set -ux \
 && apk update \
 && apk add ca-certificates git tzdata \
 && update-ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN set -ux \
 && go mod download \
 && go mod verify
COPY . .
RUN set -ux \
 && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -o mybin -ldflags '-s -w'

FROM scratch
WORKDIR /
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/mybin /mybin
LABEL org.opencontainers.image.source https://github.com/me/mybin
ENV PORT=8888
EXPOSE $PORT/tcp
ENTRYPOINT ["/mybin"]
```
outputs:
```
# locked {
# locked using:OCI-image-FROM
# locked   and:track=golang:1.16.2-alpine3.12
# locked   and:semver=1.16.2~1
FROM golang@sha256:49812b175d0a519eabcab8b70b741e3a84edb08fecfb3e7a252ed08626b13c48 AS builder
# locked }
...
```

### Bazel
Given the WORKSPACE file:
```python
# locked {
# locked using:bazel-http_archive-from
# locked   and:tag=~3
# locked   and:release=glfw-{tag_digits}.bin.MACOS.zip
# locked   and:strip_prefix=glfw-{tag_digits}.bin.MACOS
# locked   and:remote=git://github.com/glfw/glfw.git
http_archive(
    name = "glfw_osx",
    build_file = "@//third_party:glfw3_osx.BUILD",
    BREAKS "correct and usable"
)
# locked }
```
outputs:
```python
# locked {
# locked using:bazel-http_archive-from
# locked   and:tag=~3
# locked   and:release=glfw-{tag_digits}.bin.MACOS.zip
# locked   and:strip_prefix=glfw-{tag_digits}.bin.MACOS
# locked   and:remote=git://github.com/glfw/glfw.git
http_archive(
    name = "glfw_osx",
    build_file = "@//third_party:glfw3_osx.BUILD",
    sha256 = "e412c75f850c320192df491ec3bf623847fafa847b46ffd3bbd7478057148f5a",
    strip_prefix = "glfw-3.3.2.bin.MACOS",
    type = "zip",
    urls = ["https://github.com/glfw/glfw/releases/download/3.3.2/glfw-3.3.2.bin.MACOS.zip"],
)
# locked }
```
