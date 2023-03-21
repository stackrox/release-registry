// Package version provides support for semantic versions.
package version

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
)

const (
	nightlyNewReplacePart = "0-aaa"
	nightlyOldReplacePart = "x"
	nightlyIdentifier     = "-nightly-"
)

func sortVersions(versions []string) ([]string, error) {
	vs := make([]*semver.Version, len(versions))

	for i, r := range versions {
		if strings.Contains(r, nightlyIdentifier) {
			r = strings.Replace(r, nightlyOldReplacePart, nightlyNewReplacePart, 1)
		}

		v, err := semver.StrictNewVersion(r)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error parsing version: %s", v))
		}

		vs[i] = v
	}

	// Sort implicitly works for our use cases, as release candidate numbers and
	// nightly dates are lexicographically sorted.
	// For the same semVer base version (A.B.C), RCs are preferred over nightlies
	// due to lexicographic sorting.
	// Published versions are preferred over prerelease versions like RCs and nightlies.
	sort.Sort(semver.Collection(vs))

	sortedVersions := make([]string, len(versions))
	for i, r := range vs {
		sortedVersions[i] = strings.Replace(r.String(), nightlyNewReplacePart, nightlyOldReplacePart, 1)
	}

	return sortedVersions, nil
}

// LatestVersion returns the latest version from a list of unsorted semantic versions.
func LatestVersion(versions []string) (string, error) {
	vs, err := sortVersions(versions)
	if err != nil {
		return "", err
	}

	return vs[len(vs)-1], nil
}

// Validate returns true if the given version is considered valid SemVer in the context of the service.
func Validate(version string) error {
	if strings.Contains(version, nightlyIdentifier) {
		version = strings.Replace(version, nightlyOldReplacePart, nightlyNewReplacePart, 1)
	}

	_, err := semver.StrictNewVersion(version)
	if err != nil {
		return errors.Wrap(err, "version is not valid SemVer")
	}

	return nil
}
