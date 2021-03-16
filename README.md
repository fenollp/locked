# locked
Universal lockfile (not a daemon)

* `locked [--tracked=REGEX --] ( FILE )+`
* `locked .../Dockerfile`
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

### in-the-wild text
Requires a `/Lockfile`
```
Bla, bla blabla: https://github.com/fenollp/locked/blob/main/README.md
blblblbl
```
`locked .` =>
```
Lockfile
---
https://github.com/fenollp/locked/blob/main/README.md => https://github.com/fenollp/locked/blob/38004f0260b4ad77daa873b7f7428fa151771a0d/README.md
```
```
Bla, bla blabla: https://github.com/fenollp/locked/blob/38004f0260b4ad77daa873b7f7428fa151771a0d/README.md
blblblbl
```

## Lockfile
* "pattern" => "locked"
    * str => str :: mapping
        * solves the simplest case (in-the-wild text)
        * solves Dockerfile FROM
    * dict => [mapping]
        * solves complex http_archive (constraints + multiple strings outputed)
* recursive overloading
    * a la .gitignore
    * allows for a global Lockfile
* contains track=/"pattern" + method/"using" + "using"'s args
    * also contains op results
        * so they can be strstr ez

---
`../Lockfile`
```hcl
#!/usr/bin/env locked --version-is=1

# Some comment
# File is fmt'd and comments are kept when resolve gives outputs.

track "goreleaser/goreleaser" {
    using = "get-oci-image-sha256"
    gives = {  # Code generated; DO NOT EDIT
        at = "2009-11-10 23:00:00 +0000 UTC"
        tracked = "goreleaser/goreleaser@sha256:fa75344740e66e5bb55ad46426eb8e6c8dedbd3dcfa15ec1c41897b143214ae2"
    }
}
```
---
`~/.config/Lockfile`
```hcl
# Crazy: generic Docker image locking!
# Looks for images used in Dockerfiles, then updates all mentions of said images
# NOTE: images used but not found by the following regex won't be locked
track "^from\s+(?:--[^\s]+\s+)*([^\s]+)" { # Considered a regex iff starts (or ends) with ^ (or $)
                                           # MUST have exactly 1 capturing group
                                           # Multiline regexp by default
                                           # Case-sensitive iff contains uppercase (in which case (?ms) is prepended to the regex, instead of (?ims))
    using = "get-oci-image-sha256"

    # Code generated; DO NOT EDIT
    tracking "FROM --platform=$BUILDPLATFORM goreleaser/goreleaser AS go-releaser" {
        at = "2009-11-10 23:00:00 +0000 UTC"
        gives = "FROM --platform=$BUILDPLATFORM goreleaser/goreleaser@sha256:fa75344740e66e5bb55ad46426eb8e6c8dedbd3dcfa15ec1c41897b143214ae2 AS go-releaser"
    }
    tracking "from python:3.6" {
        at = "2009-11-10 23:00:00 +0000 UTC"
        gives = "from python@sha256:d264d2c9cf7d4b43908150e0bcd2eefe275aba956ff718169fc3c1f7727a0d0a"
    }
}

# https://docs.github.com/en/github/managing-files-in-a-repository/getting-permanent-links-to-files
track "^https://github.com/[^/]+/[^/]+/(?:blame|blob|tree)/([^/]+)/[^/]+" {
    using = "permalink-github"
}

# Resolves to URL to latest release
track "^(https://github.com/[^/]+/[^/]+/releases/latest" {
    using = "permalink-github-latest-release"
}
```
---
* re-write Lockfile on non-zero resolutions
* enforce hashbang (mentions version)
    * so language can be changed
    * actually don't: forward compat MUST be ensured
    * On first run: prepend UTC datetime to all Lockfile.s within git repo
    * On next runs:
        * only select predefined rules as they were at that time
        * prepend that time to all new Lockfile.s
    * On next runs and after 3 months: warn about running `--upgrade`
        * this (after checking there are no git changes) sets times to now then locks
* resolve Lockfile from XDG, then .git/Lockfile, then recursively down
    * exactly like .gitignore.s
    * inner Lockfile has priority over parent
        * error when `>1 "track" && !=1 "using"` in same Lockfile though
    * allow extensions: `(none)` | `.hcl` | `.json` | `.yaml`
        * but disallow more than one Lockfile.* per directory
* error out if . isn't within a git repo || . has changes
* `using = "get-oci-image-sha256"` is
    1. trim tracked string
    1. `docker pull '${TRACKED}' && docker inspect --format='{{.RepoDigests}}' '${TRACKED}'`
    1. ensure output is list of length 1 then resolve to first element
    1. remove image if it was we that just pulled it (good enough heuristic: worst case is image is pulled twice)
* when is TTY print output of each running resolver in parallel, like `DOCKER_BUILDKIT=1 docker build`
* opt-out a line from rewriting by adding a line just above that contains `Lockfile: skip next line`
* plugin system (based on [Dockerfile frontend syntaxes](https://github.com/moby/buildkit/blob/4eca10a46c7f309582e60dcc52b54fe7a5c7e3d2/frontend/dockerfile/docs/syntax.md))
    * `service Lockfile { rpc Track(TrackReq) returns (TrackRep); }`
        * `TrackReq { Track, Whole string }`
            * `TrackReq.Track`: the value or the capturing group
            * `TrackReq.Whole`: the value or the whole regex match (can be multiple lines)
        * `TrackRep { Gives string }`
            * `TrackRep.Gives`: locked value with `Whole` context (i.e. `TrackReq.Whole` with `TrackReq.Track` replaced)
