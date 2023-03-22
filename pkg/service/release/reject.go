package release

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
)

func (s *server) Reject(
	ctx context.Context, rejectedRelease *v1.ReleaseServiceRejectRequest,
) (*v1.ReleaseServiceRejectResponse, error) {
	tag := rejectedRelease.GetTag()
	release, err := models.RejectRelease(tag, rejectedRelease.GetPreload())

	if err != nil {
		message := "could not reject release"
		log.Infow(message, "tag", tag, "error", err.Error())

		return nil, errors.Wrap(err, fmt.Sprintf("%s '%s'", message, tag))
	}

	releaseResponse := newRejectReleaseResponseFromRelease(release)

	return releaseResponse, nil
}

func newRejectReleaseResponseFromRelease(release *models.Release) *v1.ReleaseServiceRejectResponse {
	return &v1.ReleaseServiceRejectResponse{
		Meta:     newMetaFromRelease(release),
		Tag:      release.Tag,
		Commit:   release.Commit,
		Creator:  release.Creator,
		Metadata: newReleaseMetadataFromRelease(release),
		Rejected: release.Rejected,
	}
}
