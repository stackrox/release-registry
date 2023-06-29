package release

import (
	"context"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/pkg/utils/conversions"
)

func (s *releaseImpl) Update(
	ctx context.Context, updatedRelease *v1.ReleaseServiceUpdateRequest,
) (*v1.ReleaseServiceUpdateResponse, error) {
	errMessage := "could not update release"
	tag := updatedRelease.GetTag()

	releaseMetadata := []models.ReleaseMetadata{}
	for _, metadata := range updatedRelease.GetMetadata() {
		releaseMetadata = append(releaseMetadata, models.ReleaseMetadata{
			Key:   metadata.GetKey(),
			Value: metadata.GetValue(),
		})
	}

	release, err := models.UpdateRelease(tag, releaseMetadata, updatedRelease.GetIncludeRejected())

	if err != nil {
		log.Infow(errMessage, "tag", tag, "error", err.Error())

		return nil, errors.WithMessagef(err, "%s '%s'", errMessage, tag)
	}

	updateReleaseResponse := conversions.NewUpdateReleaseResponseFromRelease(release)

	return updateReleaseResponse, nil
}
