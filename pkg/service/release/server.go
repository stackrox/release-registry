// Package release contains the server to manage Releases.
package release

import (
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	v1Meta "github.com/stackrox/release-registry/gen/go/proto/shared/v1"
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/logging"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func newReleaseMetadataFromRelease(release *models.Release) []*v1.ReleaseMetadata {
	releaseMetadata := []*v1.ReleaseMetadata{}

	for _, metadata := range release.Metadata {
		releaseMetadata = append(releaseMetadata, &v1.ReleaseMetadata{
			Key:   metadata.Key,
			Value: metadata.Value,
		})
	}

	return releaseMetadata
}

func newV1QualityMilestoneFromRelease(release *models.Release) []*v1.QualityMilestone {
	qualityMilestones := []*v1.QualityMilestone{}
	for _, qm := range release.QualityMilestones {
		qualityMilestones = append(qualityMilestones, &v1.QualityMilestone{
			Name:     qm.QualityMilestoneDefinition.Name,
			Approver: qm.Approver,
			Meta: &v1Meta.Meta{
				Id:        int64(qm.ID),
				CreatedAt: timestamppb.New(qm.CreatedAt),
				UpdatedAt: timestamppb.New(qm.UpdatedAt),
			},
		})
	}

	return qualityMilestones
}

func newMetaFromRelease(release *models.Release) *v1Meta.Meta {
	return &v1Meta.Meta{
		Id:        int64(release.ID),
		CreatedAt: timestamppb.New(release.CreatedAt),
		UpdatedAt: timestamppb.New(release.UpdatedAt),
	}
}
