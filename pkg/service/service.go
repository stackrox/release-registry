// Package service contains only a stub function.
package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/logging"
	helloworld "github.com/stackrox/release-registry/pkg/service/hello-world"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//nolint:gochecknoglobals
var log = logging.CreateProductionLogger()

// Run runs a gRPC + HTTP server on a single port.
func Run(config *configuration.Config) (<-chan error, error) {
	// listenAddress is the address that the server will listen on.
	listenAddress := fmt.Sprintf("0.0.0.0:%d", config.Server.Port)
	// connectAddress is the address that the gateway client will connect to.
	connectAddress := fmt.Sprintf("localhost:%d", config.Server.Port)

	mux := http.NewServeMux()
	errCh := make(chan error, 1)

	server := grpc.NewServer()

	// TODO: loop here over all APIs
	v1.RegisterHelloWorldServiceServer(server, helloworld.NewHelloWorldServer())

	// muxHandler is a HTTP handler that can route both HTTP/2 gRPC and HTTP1.1
	// requests.
	muxHandler := http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		// Is the current request gRPC?
		if request.ProtoMajor == 2 && strings.HasPrefix(request.Header.Get("Content-Type"), "application/grpc") {
			server.ServeHTTP(responseWriter, request)

			return
		}

		// Fallback to HTTP
		mux.ServeHTTP(responseWriter, request)
	})

	log.Infow("starting gRPC server", "listenAddress", listenAddress)

	go func() {
		// TODO: make TLS
		if err := http.ListenAndServe(listenAddress, h2c.NewHandler(muxHandler, &http2.Server{})); err != nil {
			errCh <- err
		}
	}()

	dialOptions := []grpc.DialOption{
		grpc.WithBlock(),
		// TODO: make secure credentials
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	log.Infow("starting gRPC-Gateway client", "connectAddress", connectAddress)
	conn, err := grpc.Dial(connectAddress, dialOptions...)

	if err != nil {
		return nil, errors.Wrap(err, "dialing gRPC")
	}

	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption("*", &runtime.JSONPb{}),
	)

	// Register each service
	// TODO: make loop
	if err = v1.RegisterHelloWorldServiceHandler(context.Background(), gwMux, conn); err != nil {
		return nil, errors.Wrap(err, "error registering hello world service handler")
	}

	routeMux := http.NewServeMux()
	routeMux.Handle("/", gwMux)

	mux.Handle("/", routeMux)

	return errCh, nil
}
