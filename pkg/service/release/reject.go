package release

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/pkg/utils/conversions"
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

	releaseResponse := conversions.NewRejectReleaseResponseFromRelease(release)

	return releaseResponse, nil
}
