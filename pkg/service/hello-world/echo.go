// Package helloworld contains the echo server to demonstrate the setup.
package helloworld

import (
	"context"
	"strings"

	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/logging"
)

//nolint:gochecknoglobals
var log = logging.CreateProductionLogger()

type server struct {
	v1.UnimplementedHelloWorldServiceServer
}

// NewHelloWorldServer registers a new grpc server for this package.
func NewHelloWorldServer() *server {
	return &server{}
}

func (*server) Echo(ctx context.Context, message *v1.EchoRequest) (*v1.EchoResponse, error) {
	log.Infow("received message", "message", message.GetValue())
	response := &v1.EchoResponse{
		Value: strings.ToUpper(message.GetValue()),
	}

	return response, nil
}
