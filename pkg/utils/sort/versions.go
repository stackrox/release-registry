// Package sort provides order for semantic versions.
package sort

import (
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
)

func sortVersions(versions []string) ([]*semver.Version, error) {
	vs := make([]*semver.Version, len(versions))

	for i, r := range versions {
		v, err := semver.NewVersion(r)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error parsing version: %s", v))
		}

		vs[i] = v
	}

	sort.Sort(semver.Collection(vs))

	return vs, nil
}

// LatestVersion returns the latest version from a list of unsorted semantic versions.
func LatestVersion(versions []string) (string, error) {
	vs, err := sortVersions(versions)
	if err != nil {
		return "", err
	}

	return vs[len(vs)-1].String(), nil
}
