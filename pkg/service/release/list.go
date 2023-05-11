package release

import (
	"context"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/pkg/utils/conversions"
)

//nolint:cyclop
func (s *releaseImpl) List(
	ctx context.Context, listRelease *v1.ReleaseServiceListRequest,
) (*v1.ReleaseServiceListResponse, error) {
	var (
		releases []models.Release
		err      error
	)

	switch {
	case listRelease.Prefix == "" && listRelease.QualityMilestoneName == "":
		releases, err = models.ListAllReleases(
			determineIgnoredReleaseKinds(listRelease),
			listRelease.GetPreload(),
			listRelease.GetIncludeRejected(),
		)
	case listRelease.Prefix == "" && listRelease.QualityMilestoneName != "":
		releases, err = models.ListAllReleasesAtQualityMilestone(
			listRelease.GetQualityMilestoneName(),
			determineIgnoredReleaseKinds(listRelease),
			listRelease.GetPreload(),
			listRelease.GetIncludeRejected(),
		)
	case listRelease.Prefix != "" && listRelease.QualityMilestoneName == "":
		releases, err = models.ListAllReleasesWithPrefix(
			listRelease.GetPrefix(),
			determineIgnoredReleaseKinds(listRelease),
			listRelease.GetPreload(),
			listRelease.GetIncludeRejected(),
		)
	case listRelease.Prefix != "" && listRelease.QualityMilestoneName != "":
		releases, err = models.ListAllReleasesWithPrefixAtQualityMilestone(
			listRelease.GetPrefix(),
			listRelease.GetQualityMilestoneName(),
			determineIgnoredReleaseKinds(listRelease),
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
