# Server
FROM golang:1.19.5 as server

WORKDIR /go/src/github.com/stackrox/release-registry

COPY cmd cmd
COPY pkg pkg
COPY gen gen
COPY go.mod .
COPY go.sum .
COPY Makefile Makefile

RUN go build -o build/release-registry -a -ldflags '-linkmode external -extldflags "-static"' cmd/server/main.go

# UI
FROM node:18.15.0 as ui

COPY ui ui

WORKDIR /ui

RUN --mount=type=secret,id=npmrc,target=/root/.npmrc make

# Application
FROM scratch as app

COPY --from=server /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=server /go/src/github.com/stackrox/release-registry/build/release-registry /release-registry

COPY --from=server /go/src/github.com/stackrox/release-registry/gen/openapiv2 docs

COPY --from=ui ui/build ui

ENTRYPOINT [ "/release-registry" ]
