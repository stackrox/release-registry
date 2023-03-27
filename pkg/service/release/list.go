package release

import (
	"context"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/pkg/utils/conversions"
)

//nolint:cyclop
func (s *server) List(
	ctx context.Context, listRelease *v1.ReleaseServiceListRequest,
) (*v1.ReleaseServiceListResponse, error) {
	var (
		releases []models.Release
		err      error
	)

	switch {
	case listRelease.Prefix == nil && listRelease.QualityMilestoneName == nil:
		releases, err = models.ListAllReleases(
			listRelease.GetPreload(),
			listRelease.GetIncludeRejected(),
		)
	case listRelease.Prefix == nil && listRelease.QualityMilestoneName != nil:
		releases, err = models.ListAllReleasesAtQualityMilestone(
			listRelease.GetQualityMilestoneName(),
			listRelease.GetPreload(),
			listRelease.GetIncludeRejected(),
		)
	case listRelease.Prefix != nil && listRelease.QualityMilestoneName == nil:
		releases, err = models.ListAllReleasesWithPrefix(
			listRelease.GetPrefix(),
			listRelease.GetPreload(),
			listRelease.GetIncludeRejected(),
		)
	case listRelease.Prefix != nil && listRelease.QualityMilestoneName != nil:
		releases, err = models.ListAllReleasesWithPrefixAtQualityMilestone(
			listRelease.GetPrefix(),
			listRelease.GetQualityMilestoneName(),
			listRelease.GetPreload(),
			listRelease.GetIncludeRejected(),
		)
	default:
		err = errors.New("unexpected combination of prefix and quality milestone name set")
	}

	if err != nil {
		message := "could not list releases"
		log.Infow(message, "error", err.Error())

		return nil, errors.Wrap(err, message)
	}

	releaseResponse := conversions.NewListReleaseResponseFromReleases(releases)

	return releaseResponse, nil
}
