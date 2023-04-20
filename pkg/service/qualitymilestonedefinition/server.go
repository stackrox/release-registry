// Package qualitymilestonedefinition contains the server to manage QualityMilestoneDefinitions.
package qualitymilestonedefinition

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/stackrox/infra-auth-lib/service/middleware"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/logging"
	"google.golang.org/grpc"
)

type qualityMilestoneDefinitionImpl struct {
	v1.UnimplementedQualityMilestoneDefinitionServiceServer
}

var (
	_ middleware.APIService                      = (*qualityMilestoneDefinitionImpl)(nil)
	_ v1.QualityMilestoneDefinitionServiceServer = (*qualityMilestoneDefinitionImpl)(nil)
)

//nolint:gochecknoglobals
var log = logging.CreateProductionLogger()

// NewQualityMilestoneDefinitionService creates a new QualityMilestoneDefinitionService.
//
//nolint:ireturn,nolintlint
func NewQualityMilestoneDefinitionService() (middleware.APIService, error) {
	return &qualityMilestoneDefinitionImpl{}, nil
}

// RegisterServiceServer registers this service with the given gRPC Server.
func (s *qualityMilestoneDefinitionImpl) RegisterServiceServer(server *grpc.Server) {
	v1.RegisterQualityMilestoneDefinitionServiceServer(server, s)
}

// RegisterServiceHandler registers this service with the given gRPC Gateway endpoint.
func (s *qualityMilestoneDefinitionImpl) RegisterServiceHandler(
	ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn,
) error {
	if err := v1.RegisterQualityMilestoneDefinitionServiceHandler(ctx, mux, conn); err != nil {
		return errors.Wrap(err, "could not register qualitymilestonedefinition service handler")
	}

	return nil
}

// Access configures access for this service.
func (s *qualityMilestoneDefinitionImpl) Access() map[string]middleware.Access {
	return map[string]middleware.Access{
		"/proto.api.v1.QualityMilestoneDefinitionService/Create": middleware.Authenticated,
		"/proto.api.v1.QualityMilestoneDefinitionService/Get":    middleware.Anonymous,
		"/proto.api.v1.QualityMilestoneDefinitionService/List":   middleware.Anonymous,
	}
}
