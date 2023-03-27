package conversions

import (
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	v1Meta "github.com/stackrox/release-registry/gen/go/proto/shared/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewQualityMilestoneMetadataFromV1QualityMilestoneMetadata converts QualityMilestoneMetadata
// from service to storage representation.
func NewQualityMilestoneMetadataFromV1QualityMilestoneMetadata(
	v1Metadata []*v1.QualityMilestoneMetadata,
) []models.QualityMilestoneMetadata {
	qualityMilestoneMetadata := make([]models.QualityMilestoneMetadata, len(v1Metadata))

	for i := range v1Metadata {
		qualityMilestoneMetadata[i] = models.QualityMilestoneMetadata{
			Key:   v1Metadata[i].GetKey(),
			Value: v1Metadata[i].GetValue(),
		}
	}

	return qualityMilestoneMetadata
}

// newV1QualityMilestonesFromQualityMilestones converts QualityMilestones from service to storage representation.
func newV1QualityMilestonesFromQualityMilestones(qualityMilestones []models.QualityMilestone) []*v1.QualityMilestone {
	v1QualityMilestones := make([]*v1.QualityMilestone, len(qualityMilestones))
	for i, qm := range qualityMilestones {
		v1QualityMilestones[i] = &v1.QualityMilestone{
			Name:     qm.QualityMilestoneDefinition.Name,
			Approver: qm.Approver,
			Meta: &v1Meta.Meta{
				Id:        int64(qm.ID),
				CreatedAt: timestamppb.New(qm.CreatedAt),
				UpdatedAt: timestamppb.New(qm.UpdatedAt),
			},
		}
	}

	return v1QualityMilestones
}

// newV1QualityMilestoneMetadataFromQualityMilestoneMetadata converts
// QualityMilestoneMetadata from service to storage representation.
func newV1QualityMilestoneMetadataFromQualityMilestoneMetadata(
	metadata []models.QualityMilestoneMetadata,
) []*v1.QualityMilestoneMetadata {
	qualityMilestoneMetadata := make([]*v1.QualityMilestoneMetadata, len(metadata))
	for i, m := range metadata {
		qualityMilestoneMetadata[i] = &v1.QualityMilestoneMetadata{
			Key:   m.Key,
			Value: m.Value,
		}
	}

	return qualityMilestoneMetadata
}

// newQualityMilestonesFromV1QualityMilestones converts QualityMilestones from service to storage representation.
func newQualityMilestonesFromV1QualityMilestones(qmList []*v1.QualityMilestone) []models.QualityMilestone {
	qualityMilestones := make([]models.QualityMilestone, len(qmList))
	for i := range qmList {
		qualityMilestones[i] = models.QualityMilestone{
			Approver: qmList[i].GetApprover(),
			Metadata: NewQualityMilestoneMetadataFromV1QualityMilestoneMetadata(qmList[i].GetMetadata()),
		}
	}

	return qualityMilestones
}
