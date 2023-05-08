package version_test

import (
	"testing"

	"github.com/stackrox/release-registry/pkg/utils/version"
	"github.com/stretchr/testify/assert"
)

func TestLatestDev(t *testing.T) {
	assert.Equal(t, version.GetKind("1.2.x-24-gabcdef1234"), version.DevelopmentKind)
	latest, err := version.LatestVersion([]string{"1.2.x-24-gabcdef1234", "1.2.x-123-gabcdef1234"})
	assert.NoError(t, err)
	assert.Equal(t, latest, "1.2.3-123-gabcdef1234")
}

func TestLatestReleaseVersion(t *testing.T) {
	latest, err := version.LatestVersion([]string{"1.2.3", "2.3.4", "1.2.4"})
	assert.NoError(t, err)
	assert.Equal(t, latest, "2.3.4")
}

func TestLatestVersionRC(t *testing.T) {
	latest, err := version.LatestVersion([]string{"1.2.2", "1.2.3-rc.1", "1.2.3-rc.10", "1.2.3-rc.2", "1.2.2-rc.2"})
	assert.NoError(t, err)
	assert.Equal(t, "1.2.3-rc.10", latest)
}

func TestLatestVersionNightly(t *testing.T) {
	latest, err := version.LatestVersion([]string{
		"1.2.x-nightly-20230320",
		"1.2.x-nightly-20230212",
		"1.2.x-nightly-20230321",
		"1.2.x-nightly-20230319",
	})
	assert.NoError(t, err)
	assert.Equal(t, "1.2.x-nightly-20230321", latest)
}

func TestLatestVersionMixed(t *testing.T) {
	latest, err := version.LatestVersion([]string{
		"1.2.x-nightly-20230320",
		"1.2.x-nightly-20230212",
		"1.2.x-nightly-20230321",
		"1.2.x-nightly-20230319",
		"1.2.3",
		"1.2.3-rc.1",
		"1.2.3-rc.10",
		"1.2.3-rc.2",
		"1.2.2-rc.2",
	})

	assert.NoError(t, err)
	assert.Equal(t, "1.2.x-nightly-20230321", latest)
}
