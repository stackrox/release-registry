// Package qualitymilestonedefinition contains the server to manage QualityMilestoneDefinitions.
package qualitymilestonedefinition

import (
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/logging"
)

//nolint:gochecknoglobals
var log = logging.CreateProductionLogger()

type server struct {
	v1.UnimplementedQualityMilestoneDefinitionServiceServer

	Config *configuration.Config
}

// NewQualityMilestoneDefinitionServer registers a new grpc server for this package.
func NewQualityMilestoneDefinitionServer(config *configuration.Config) *server {
	return &server{Config: config}
}
