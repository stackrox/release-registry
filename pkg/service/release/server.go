// Package release contains the server to manage Releases.
package release

import (
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/logging"
)

//nolint:gochecknoglobals
var log = logging.CreateProductionLogger()

type server struct {
	v1.UnimplementedReleaseServiceServer

	Config *configuration.Config
}

// NewReleaseServer registers a new grpc server for this package.
func NewReleaseServer(config *configuration.Config) *server {
	return &server{Config: config}
}
