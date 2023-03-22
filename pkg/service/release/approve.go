package release

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	v1Meta "github.com/stackrox/release-registry/gen/go/proto/shared/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *server) Approve(
	ctx context.Context, approvedRelease *v1.ReleaseServiceApproveRequest,
) (*v1.ReleaseServiceApproveResponse, error) {
	tag := approvedRelease.GetTag()
	qualityMilestoneName := approvedRelease.GetQualityMilestoneDefinitionName()
	approver := approvedRelease.GetApprover()

	qualityMilestoneMetadata := newQualityMilestoneMetadataFromV1QualityMilestoneMetadata(approvedRelease.GetMetadata())

	qualityMilestone, err := models.ApproveQualityMilestone(
		s.Config, tag, qualityMilestoneName, approver, qualityMilestoneMetadata,
	)
	if err != nil {
		message := "could not approve quality milestone for release"
		log.Infow(message, "tag", tag, "qualityMilestone", qualityMilestoneName, "approver", approver, "error", err.Error())

		return nil, errors.Wrap(
			err,
			fmt.Sprintf(
				"%s (Release %s, Quality Milestone %s, Approver %s)",
				message, tag, qualityMilestoneName, approver,
			),
		)
	}

	response := newApproveQualityMilestoneResponseFromQualityMilestone(qualityMilestone)

	return response, nil
}

func newApproveQualityMilestoneResponseFromQualityMilestone(
	qm *models.QualityMilestone,
) *v1.ReleaseServiceApproveResponse {
	return &v1.ReleaseServiceApproveResponse{
		Meta:                           newMetaFromQualityMilestone(qm),
		Tag:                            qm.Release.Tag,
		QualityMilestoneDefinitionName: qm.QualityMilestoneDefinition.Name,
		Approver:                       qm.Approver,
		Metadata:                       newV1QualityMilestoneMetadataFromModelQualityMilestoneMetadata(qm.Metadata),
	}
}

func newMetaFromQualityMilestone(qm *models.QualityMilestone) *v1Meta.Meta {
	return &v1Meta.Meta{
		Id:        int64(qm.ID),
		CreatedAt: timestamppb.New(qm.CreatedAt),
		UpdatedAt: timestamppb.New(qm.UpdatedAt),
	}
}

func newV1QualityMilestoneMetadataFromModelQualityMilestoneMetadata(
	metadata []models.QualityMilestoneMetadata,
) []*v1.QualityMilestoneMetadata {
	qualityMilestoneMetadata := make([]*v1.QualityMilestoneMetadata, len(metadata))
	for i := range metadata {
		qualityMilestoneMetadata[i] = &v1.QualityMilestoneMetadata{
			Key:   metadata[i].Key,
			Value: metadata[i].Value,
		}
	}

	return qualityMilestoneMetadata
}

func newQualityMilestoneMetadataFromV1QualityMilestoneMetadata(
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
