// Package version provides support for semantic versions.
package version

import (
	"sort"

	"github.com/pkg/errors"
)

// LatestVersion returns the latest version from a list of unsorted semantic versions.
func LatestVersion(versions []string) (string, error) {
	sort.Sort(ByVersion(versions))

	return versions[len(versions)-1], nil
}

// Validate returns no error if the given version is a valid version in the context of the service.
func Validate(version string) error {
	kind := GetKind(version)
	if kind == InvalidKind {
		return errors.Errorf("%s is not a valid version", version)
	}

	return nil
}
