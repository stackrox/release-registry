package conversions

import (
	v1Meta "github.com/stackrox/release-registry/gen/go/proto/shared/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newV1MetaFromRelease(release *models.Release) *v1Meta.Meta {
	return &v1Meta.Meta{
		Id:        int64(release.ID),
		CreatedAt: timestamppb.New(release.CreatedAt),
		UpdatedAt: timestamppb.New(release.UpdatedAt),
	}
}

func newV1MetaFromQualityMilestone(qm *models.QualityMilestone) *v1Meta.Meta {
	return &v1Meta.Meta{
		Id:        int64(qm.ID),
		CreatedAt: timestamppb.New(qm.CreatedAt),
		UpdatedAt: timestamppb.New(qm.UpdatedAt),
	}
}

func newV1MetaFromQualityMilestoneDefinition(qmd *models.QualityMilestoneDefinition) *v1Meta.Meta {
	return &v1Meta.Meta{
		Id:        int64(qmd.ID),
		CreatedAt: timestamppb.New(qmd.CreatedAt),
		UpdatedAt: timestamppb.New(qmd.UpdatedAt),
	}
}
