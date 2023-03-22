package release

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
)

func (s *server) Get(
	ctx context.Context, getRelease *v1.ReleaseServiceGetRequest,
) (*v1.ReleaseServiceGetResponse, error) {
	tag := getRelease.GetTag()
	release, err := models.GetRelease(
		tag,
		getRelease.GetPreload(),
		getRelease.GetIncludeRejected(),
	)

	if err != nil {
		message := "could not get release"
		log.Infow(message, "tag", tag, "error", err.Error())

		return nil, errors.Wrap(err, fmt.Sprintf("%s '%s'", message, tag))
	}

	releaseResponse := newGetReleaseResponseFromRelease(release)

	return releaseResponse, nil
}

func newGetReleaseResponseFromRelease(release *models.Release) *v1.ReleaseServiceGetResponse {
	return &v1.ReleaseServiceGetResponse{
		Meta:              newMetaFromRelease(release),
		Tag:               release.Tag,
		Commit:            release.Commit,
		Creator:           release.Creator,
		Metadata:          newReleaseMetadataFromRelease(release),
		Rejected:          release.Rejected,
		QualityMilestones: newV1QualityMilestoneFromRelease(release),
	}
}
