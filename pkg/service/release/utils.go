package release

import (
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/utils/version"
)

type releaseRequestWithIgnoredKinds interface {
	GetIgnoredReleaseKinds() []v1.ReleaseKind
}

func determineIgnoredReleaseKinds[R releaseRequestWithIgnoredKinds](release R) []version.Kind {
	ignoredKinds := []version.Kind{}
	releaseKinds := release.GetIgnoredReleaseKinds()

	for _, kind := range releaseKinds {
		ignoredKinds = append(ignoredKinds, version.Kind(kind))
	}

	return ignoredKinds
}
