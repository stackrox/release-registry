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

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/stackrox/infra-auth-lib/auth"
	authService "github.com/stackrox/infra-auth-lib/service"
	"github.com/stackrox/infra-auth-lib/service/middleware"
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/logging"
	"github.com/stackrox/release-registry/pkg/service/healthz"
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
	Config   *configuration.Config
	ErrCh    chan error
	Services []middleware.APIService
	oidc     auth.OidcAuth
}

//nolint:gochecknoglobals
var log = logging.CreateProductionLogger()

// New creates and returns a new server instance.
func New(config *configuration.Config, oidc auth.OidcAuth, services ...middleware.APIService) Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// Extract user from JWT token stored in HTTP cookie.
			middleware.ContextInterceptor(middleware.UserEnricher(oidc)),
			// Extract service-account from token stored in Authorization header.
			middleware.ContextInterceptor(middleware.ServiceAccountEnricher(oidc.ValidateServiceAccountToken)),

			// Enrich with admin privileges from password stored in Authorization header.
			middleware.ContextInterceptor(middleware.AdminEnricher(config.Tenant.Password)),

			// Enforce authenticated user access on resources that declare it.
			middleware.ContextInterceptor(middleware.EnforceAccess),
		)),
	)
	s := Server{
		grpcServer,
		config,
		make(chan error, 1),
		services,
		oidc,
	}

	return s
}

// Run runs a gRPC + HTTP server on a single port.
func (s *Server) Run() error {
	mux := s.initMux()

	s.registerServiceServers()

	gwMux, err := s.initGateway()
	if err != nil {
		return err
	}

	mux.Handle("/", s.serveApplicationResources())
	mux.Handle("/healthz/", healthz.NewHandler())
	mux.Handle("/docs/", s.newDocsHandler())
	mux.Handle("/v1/", gwMux)

	s.oidc.Handle(mux)

	return nil
}

func (s *Server) registerServiceServers() {
	// Register the gRPC API service.
	for _, apiSvc := range s.Services {
		apiSvc.RegisterServiceServer(s.Server)
	}
}

func (s *Server) registerServiceHandlers(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	for _, apiSvc := range s.Services {
		if err := apiSvc.RegisterServiceHandler(context.Background(), mux, conn); err != nil {
			return errors.Wrap(err, "cannot register service handler")
		}
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
	if err = s.registerServiceHandlers(gwMux, conn); err != nil {
		return nil, err
	}

	return gwMux, nil
}

func (s *Server) newDocsHandler() http.Handler {
	return http.StripPrefix(
		"/docs/",
		http.FileServer(http.Dir(s.Config.Server.DocsDir)),
	)
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

// CreateAPIServices returns a list of all gRPC API services.
func CreateAPIServices(config *configuration.Config, oidc auth.OidcAuth) ([]middleware.APIService, error) {
	//nolint:wrapcheck
	return middleware.Services(
		func() (middleware.APIService, error) {
			//nolint:wrapcheck
			return release.NewReleaseService(config.Tenant.EmailDomain)
		},
		qualitymilestonedefinition.NewQualityMilestoneDefinitionService,
		func() (middleware.APIService, error) {
			//nolint:wrapcheck
			return authService.NewUserService(oidc.GenerateServiceAccountToken)
		},
	)
}

// serveApplicationResources handles requests for SPA endpoints as well as
// regular resources.
func (s *Server) serveApplicationResources() http.Handler {
	type rule struct {
		path      string
		spa       bool
		anonymous bool
		prefix    bool
	}

	// List of path rules, roughly ordered from most-likely matched to
	// least-likely matched.
	rules := []rule{
		{path: "/static/", prefix: true},
		{path: "/manifest.json"},
		{path: "/favicon.ico", anonymous: true},
		{path: "/logout-page.html", anonymous: true},
		{path: "/", spa: true, prefix: true},
	}

	fs := http.FileServer(http.Dir(s.Config.Server.StaticDir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path

		for _, rule := range rules {
			if rule.prefix {
				// If this rule is supposed to match a path prefix, and that
				// path prefix isn't matched, move onto the next rule.
				if !strings.HasPrefix(requestPath, rule.path) {
					continue
				}
			} else {
				// If this rule is supposed to match a path exactly, and that
				// path isn't exactly matched, move onto the next rule.
				if requestPath != rule.path {
					continue
				}
			}

			// If the path is a path in the SPA, set the path to be the root,
			// so that the index.html is served.
			if rule.spa {
				r.URL.Path = "/"
			}

			if rule.anonymous {
				// Serve this path anonymously (without any authentication).
				fs.ServeHTTP(w, r)
			} else {
				// Serve this path with authentication.
				s.oidc.Authorized(fs).ServeHTTP(w, r)
			}

			return
		}
		// No rules matched, so serve this path with authentication by default.
		s.oidc.Authorized(fs).ServeHTTP(w, r)
	})
}
