# Builder
FROM golang:1.19.5 as builder

WORKDIR /go/src/github.com/stackrox/release-registry

COPY cmd cmd
COPY pkg pkg
COPY gen gen
COPY go.mod .
COPY go.sum .
COPY Makefile Makefile

RUN go build -o build/release-registry -a -ldflags '-linkmode external -extldflags "-static"' cmd/server/main.go

# Application
FROM scratch as app

COPY --from=builder /go/src/github.com/stackrox/release-registry/build/release-registry /release-registry

COPY --from=builder /go/src/github.com/stackrox/release-registry/gen/openapiv2/proto gen/openapiv2/proto

ENTRYPOINT [ "/release-registry" ]
