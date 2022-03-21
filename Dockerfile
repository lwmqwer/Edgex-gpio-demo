#
# Copyright (c) 2022 Wuming Liu
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.17-alpine AS builder
ARG TARGETOS TARGETARCH
WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY main.go main.go
COPY version.go version.go
COPY driver/ driver/

# Build
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -o device-SAK main.go version.go

# Next image - Copy built Go binary into new workspace
FROM --platform=$TARGETPLATFORM alpine:3.14

ENV APP_PORT=41000
#expose command data port
EXPOSE $APP_PORT

WORKDIR /
COPY --from=builder /workspace/device-SAK .
COPY /res/ /

ENTRYPOINT ["/device-SAK"]
CMD ["-cp=consul.http://edgex-core-consul:8500", "--registry", "--confdir=/res"]
