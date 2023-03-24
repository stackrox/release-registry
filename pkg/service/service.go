// Package service contains only a stub function.
package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/logging"
	"github.com/stackrox/release-registry/pkg/service/qualitymilestonedefinition"
	"github.com/stackrox/release-registry/pkg/service/release"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Server is a wrapped gRPC server with config and error channel.
type Server struct {
	*grpc.Server
	Config *configuration.Config
	ErrCh  chan error
}

//nolint:gochecknoglobals
var log = logging.CreateProductionLogger()

// New creates and returns a new server instance.
func New(config *configuration.Config) Server {
	s := Server{
		grpc.NewServer(),
		config,
		make(chan error, 1),
	}
	s.registerServiceServers()

	return s
}

func (s *Server) registerServiceServers() {
	v1.RegisterQualityMilestoneDefinitionServiceServer(
		s,
		qualitymilestonedefinition.NewQualityMilestoneDefinitionServer(s.Config),
	)
	v1.RegisterReleaseServiceServer(s, release.NewReleaseServer(s.Config))
}

func registerServiceHandlers(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	if err := v1.RegisterQualityMilestoneDefinitionServiceHandler(context.Background(), mux, conn); err != nil {
		return errors.Wrap(err, "could not register QualityMilestoneDefinition service handler")
	}

	if err := v1.RegisterReleaseServiceHandler(context.Background(), mux, conn); err != nil {
		return errors.Wrap(err, "could not register Release service handler")
	}

	return nil
}

func (s *Server) initMux() *http.ServeMux {
	// listenAddress is the address that the server will listen on.
	listenAddress := fmt.Sprintf("0.0.0.0:%d", s.Config.Server.Port)
	mux := http.NewServeMux()
	// muxHandler is a HTTP handler that can route both HTTP/2 gRPC and HTTP1.1 requests.
	muxHandler := http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		// Is the current request gRPC?
		if request.ProtoMajor == 2 && strings.HasPrefix(request.Header.Get("Content-Type"), "application/grpc") {
			s.ServeHTTP(responseWriter, request)

			return
		}

		// Fallback to HTTP
		mux.ServeHTTP(responseWriter, request)
	})

	log.Infow("starting gRPC server", "listenAddress", listenAddress)

	go func() {
		err := http.ListenAndServeTLS(
			listenAddress,
			s.Config.Server.Cert, s.Config.Server.Key,
			h2c.NewHandler(muxHandler, &http2.Server{}),
		)
		if err != nil {
			s.ErrCh <- err
		}
	}()

	return mux
}

func (s *Server) initGateway() (*runtime.ServeMux, error) {
	certOpt, err := grpcLocalCredentials(s.Config.Server.Cert)
	if err != nil {
		return nil, err
	}

	dialOptions := []grpc.DialOption{
		grpc.WithBlock(),
		certOpt,
	}

	// connectAddress is the address that the gateway client will connect to.
	connectAddress := fmt.Sprintf("localhost:%d", s.Config.Server.Port)

	log.Infow("starting gRPC-Gateway client", "connectAddress", connectAddress)

	conn, err := grpc.Dial(connectAddress, dialOptions...)
	if err != nil {
		return nil, errors.Wrap(err, "dialing grpc")
	}

	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption("*", &runtime.JSONPb{}),
	)

	// Register each service
	if err = registerServiceHandlers(gwMux, conn); err != nil {
		return nil, err
	}

	return gwMux, nil
}

func newDocsHandler() http.Handler {
	generatedDocsDir := "./gen/openapiv2/proto/api/v1"

	return http.StripPrefix(
		"/docs/",
		http.FileServer(http.Dir(generatedDocsDir)),
	)
}

// Run runs a gRPC + HTTP server on a single port.
func (s *Server) Run() error {
	mux := s.initMux()

	gwMux, err := s.initGateway()
	if err != nil {
		return err
	}

	mux.Handle("/docs/", newDocsHandler())
	mux.Handle("/v1/", gwMux)

	return nil
}

//nolint:ireturn,nolintlint
func grpcLocalCredentials(certFile string) (grpc.DialOption, error) {
	// Read the x509 PEM public certificate file
	pem, err := os.ReadFile(certFile)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read certificate file")
	}

	// Create an empty certificate pool, and add our single "CA" certificate to
	// it. This allows us to trust the local server specifically, as its
	// serving the same exact certificate.
	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("no root CA certs parsed from file %s", certFile)
	}

	return grpc.WithTransportCredentials(
		credentials.NewTLS(&tls.Config{
			RootCAs:    rootCAs,
			ServerName: "localhost",
		})), nil
}
