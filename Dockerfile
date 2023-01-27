FROM golang:1.19.5-alpine3.17 as builder

WORKDIR /usr/src/app

RUN apk add --no-cache build-base

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o build/relreg-server cmd/server/main.go

FROM alpine:3.17 as server

COPY --from=builder /usr/src/app/build/relreg-server /relreg-server

RUN mkdir /data

ENTRYPOINT [ "/relreg-server" ]
