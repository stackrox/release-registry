// Package release contains the server to manage Releases.
package release

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/stackrox/infra-auth-lib/service/middleware"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/logging"
	"google.golang.org/grpc"
)

type releaseImpl struct {
	v1.UnimplementedReleaseServiceServer

	validActorDomain string
}

var (
	_ middleware.APIService   = (*releaseImpl)(nil)
	_ v1.ReleaseServiceServer = (*releaseImpl)(nil)
)

//nolint:gochecknoglobals
var log = logging.CreateProductionLogger()

// NewReleaseService creates a new ReleaseService.
//
//nolint:ireturn,nolintlint
func NewReleaseService(validActorDomain string) (middleware.APIService, error) {
	return &releaseImpl{validActorDomain: validActorDomain}, nil
}

// RegisterServiceServer registers this service with the given gRPC Server.
func (s *releaseImpl) RegisterServiceServer(server *grpc.Server) {
	v1.RegisterReleaseServiceServer(server, s)
}

// RegisterServiceHandler registers this service with the given gRPC Gateway endpoint.
func (s *releaseImpl) RegisterServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	if err := v1.RegisterReleaseServiceHandler(ctx, mux, conn); err != nil {
		return errors.Wrap(err, "could not register release service handler")
	}

	return nil
}

// Access configures access for this service.
func (s *releaseImpl) Access() map[string]middleware.Access {
	return map[string]middleware.Access{
		"/proto.api.v1.ReleaseService/Approve":    middleware.Authenticated,
		"/proto.api.v1.ReleaseService/Create":     middleware.Authenticated,
		"/proto.api.v1.ReleaseService/FindLatest": middleware.Anonymous,
		"/proto.api.v1.ReleaseService/Get":        middleware.Anonymous,
		"/proto.api.v1.ReleaseService/List":       middleware.Anonymous,
		"/proto.api.v1.ReleaseService/Reject":     middleware.Authenticated,
		"/proto.api.v1.ReleaseService/Update":     middleware.Authenticated,
	}
}

func getActorFromContext(ctx context.Context) (string, error) {
	svcAcct, found := middleware.ServiceAccountFromContext(ctx)
	if !found {
		return "", errors.New("could not extract service account from request")
	}

	return svcAcct.Email, nil
}
