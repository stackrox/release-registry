package release

import (
	"context"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/pkg/utils/conversions"
)

//nolint:cyclop
func (s *releaseImpl) FindLatest(
	ctx context.Context, findRelease *v1.ReleaseServiceFindLatestRequest,
) (*v1.ReleaseServiceFindLatestResponse, error) {
	var (
		release *models.Release
		err     error
	)

	switch {
	case findRelease.Prefix == "" && findRelease.QualityMilestoneName == "":
		release, err = models.FindLatestRelease(
			determineIgnoredReleaseKinds(findRelease),
			findRelease.GetPreload(),
			findRelease.GetIncludeRejected(),
		)
	case findRelease.Prefix == "" && findRelease.QualityMilestoneName != "":
		release, err = models.FindLatestReleaseAtQualityMilestone(
			findRelease.GetQualityMilestoneName(),
			determineIgnoredReleaseKinds(findRelease),
			findRelease.GetPreload(),
			findRelease.GetIncludeRejected(),
		)
	case findRelease.Prefix != "" && findRelease.QualityMilestoneName == "":
		release, err = models.FindLatestReleaseWithPrefix(
			findRelease.GetPrefix(),
			determineIgnoredReleaseKinds(findRelease),
			findRelease.GetPreload(),
			findRelease.GetIncludeRejected(),
		)
	case findRelease.Prefix != "" && findRelease.QualityMilestoneName != "":
		release, err = models.FindLatestRelaseWithPrefixAtQualityMilestone(
			findRelease.GetPrefix(),
			findRelease.GetQualityMilestoneName(),
			determineIgnoredReleaseKinds(findRelease),
			findRelease.GetPreload(),
			findRelease.GetIncludeRejected(),
		)
	default:
		err = errors.New("unexpected combination of prefix and quality milestone name set")
	}

	if err != nil {
		message := "could not find latest release"
		log.Infow(message, "error", err.Error())

		return nil, errors.Wrap(err, message)
	}

	releaseResponse := conversions.NewFindLatestReleaseResponseFromRelease(release)

	return releaseResponse, nil
}
