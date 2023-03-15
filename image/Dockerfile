# Builder
FROM golang:1.19.5 as builder

ENV CGO_ENABLED=0
WORKDIR /go/src/github.com/stackrox/release-registry

COPY cmd cmd
COPY pkg pkg
COPY gen gen
COPY go.mod .
COPY go.sum .
COPY Makefile Makefile

RUN make server-binary

# Application
FROM scratch as app

COPY --from=builder /go/src/github.com/stackrox/release-registry/build/release-registry /release-registry

ENTRYPOINT [ "/release-registry" ]
