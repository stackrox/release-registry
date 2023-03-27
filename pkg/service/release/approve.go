package release

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/pkg/utils/conversions"
)

func (s *server) Approve(
	ctx context.Context, approvedRelease *v1.ReleaseServiceApproveRequest,
) (*v1.ReleaseServiceApproveResponse, error) {
	tag := approvedRelease.GetTag()
	qualityMilestoneName := approvedRelease.GetQualityMilestoneDefinitionName()
	approver := approvedRelease.GetApprover()
	qualityMilestoneMetadata := conversions.NewQualityMilestoneMetadataFromV1QualityMilestoneMetadata(
		approvedRelease.GetMetadata(),
	)

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

	response := conversions.NewApproveQualityMilestoneResponseFromQualityMilestone(qualityMilestone)

	return response, nil
}
