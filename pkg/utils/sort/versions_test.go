package sort_test

import (
	"testing"

	"github.com/stackrox/release-registry/pkg/utils/sort"
	"github.com/stretchr/testify/assert"
)

func TestLatestVersion(t *testing.T) {
	latest, err := sort.LatestVersion([]string{"1.2.3", "2.3.4", "1.2.4"})
	assert.NoError(t, err)
	assert.Equal(t, latest, "2.3.4")
}
