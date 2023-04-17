package release

import (
	"context"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/pkg/utils/conversions"
)

func (s *releaseImpl) Approve(
	ctx context.Context, approvedRelease *v1.ReleaseServiceApproveRequest,
) (*v1.ReleaseServiceApproveResponse, error) {
	errMessage := "could not approve quality milestone for release"

	approver, err := getActorFromContext(ctx)
	if err != nil {
		log.Infow(errMessage, "error", err.Error())

		return nil, errors.Wrap(err, err.Error())
	}

	tag := approvedRelease.GetTag()
	qualityMilestoneName := approvedRelease.GetQualityMilestoneDefinitionName()
	qualityMilestoneMetadata := conversions.NewQualityMilestoneMetadataFromV1QualityMilestoneMetadata(
		approvedRelease.GetMetadata(),
	)

	qualityMilestone, err := models.ApproveQualityMilestone(
		s.validActorDomain, tag, qualityMilestoneName, approver, qualityMilestoneMetadata,
	)
	if err != nil {
		log.Infow(errMessage,
			"tag", tag,
			"qualityMilestone", qualityMilestoneName,
			"approver", approver, "error",
			err.Error(),
		)

		return nil, errors.WithMessagef(
			err,
			"%s (Release %s, Quality Milestone %s, Approver %s)",
			errMessage, tag, qualityMilestoneName, approver,
		)
	}

	response := conversions.NewApproveQualityMilestoneResponseFromQualityMilestone(qualityMilestone)

	return response, nil
}
