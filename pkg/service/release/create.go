package release

import (
	"context"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/pkg/utils/conversions"
)

func (s *releaseImpl) Create(
	ctx context.Context, newRelease *v1.ReleaseServiceCreateRequest,
) (*v1.ReleaseServiceCreateResponse, error) {
	errMessage := "could not create release"

	creator, err := getActorFromContext(ctx)
	if err != nil {
		log.Infow(errMessage, "error", err.Error())

		return nil, errors.Wrap(err, err.Error())
	}

	tag := newRelease.GetTag()

	releaseMetadata := []models.ReleaseMetadata{}
	for _, metadata := range newRelease.GetMetadata() {
		releaseMetadata = append(releaseMetadata, models.ReleaseMetadata{
			Key:   metadata.GetKey(),
			Value: metadata.GetValue(),
		})
	}

	release, err := models.CreateRelease(
		s.validActorDomain,
		tag,
		newRelease.GetCommit(),
		creator,
		releaseMetadata,
	)

	if err != nil {
		log.Infow(errMessage, "tag", tag, "error", err.Error())

		return nil, errors.WithMessagef(err, "%s '%s'", errMessage, tag)
	}

	releaseResponse := conversions.NewCreateReleaseResponseFromRelease(release)

	return releaseResponse, nil
}
