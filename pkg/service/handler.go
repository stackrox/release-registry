// Package service contains only a stub function.
package service

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/logging"
	helloworld "github.com/stackrox/release-registry/pkg/service/hello-world"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//nolint:gochecknoglobals
var log = logging.CreateProductionLogger()

// Run is the entrypoint to start the services.
func Run(config *configuration.Config) {
	serverAddr := config.Server.Addr

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	grpcServer := grpc.NewServer()
	// Attach the Greeter service to the server

	v1.RegisterHelloWorldServiceServer(grpcServer, helloworld.NewHelloWorldServer())

	// Serve gRPC server
	log.Infoln("Serving gRPC", "addr", serverAddr)

	go func() {
		log.Fatalln(grpcServer.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		serverAddr,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()

	// Register Greeter handler
	err = v1.RegisterHelloWorldServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Infoln("Serving gRPC-Gateway", "addr", "0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
