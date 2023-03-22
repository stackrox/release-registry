package release

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
)

func (s *server) Create(
	ctx context.Context, newRelease *v1.ReleaseServiceCreateRequest,
) (*v1.ReleaseServiceCreateResponse, error) {
	tag := newRelease.GetTag()

	releaseMetadata := []models.ReleaseMetadata{}
	for _, metadata := range newRelease.GetMetadata() {
		releaseMetadata = append(releaseMetadata, models.ReleaseMetadata{
			Key:   metadata.GetKey(),
			Value: metadata.GetValue(),
		})
	}

	release, err := models.CreateRelease(
		s.Config,
		tag,
		newRelease.GetCommit(),
		newRelease.GetCreator(),
		releaseMetadata,
	)

	if err != nil {
		message := "could not create release"
		log.Infow(message, "tag", tag, "error", err.Error())

		return nil, errors.Wrap(err, fmt.Sprintf("%s '%s'", message, tag))
	}

	releaseResponse := newCreateReleaseResponseFromRelease(release)

	return releaseResponse, nil
}

func newCreateReleaseResponseFromRelease(release *models.Release) *v1.ReleaseServiceCreateResponse {
	return &v1.ReleaseServiceCreateResponse{
		Meta:     newMetaFromRelease(release),
		Tag:      release.Tag,
		Commit:   release.Commit,
		Creator:  release.Creator,
		Metadata: newReleaseMetadataFromRelease(release),
	}
}
