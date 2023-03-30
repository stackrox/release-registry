// Package utils contains shared test helpers.
package utils

import (
	"testing"

	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stretchr/testify/assert"
)

// AssertReleasesAreEqual compares two release models.
func AssertReleasesAreEqual(t *testing.T, expected, actual *models.Release, compareMeta, checkPreload bool) {
	t.Helper()
	assert.Equal(t, expected.Commit, actual.Commit)
	assert.Equal(t, expected.Tag, actual.Tag)
	assert.Equal(t, expected.Creator, actual.Creator)

	if compareMeta {
		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.CreatedAt.UnixNano(), actual.CreatedAt.UnixNano())
		assert.Equal(t, expected.UpdatedAt.UnixNano(), actual.UpdatedAt.UnixNano())
	}

	if checkPreload {
		for i := range expected.Metadata {
			assert.Equal(t, expected.Metadata[i].Key, actual.Metadata[i].Key)
			assert.Equal(t, expected.Metadata[i].Value, actual.Metadata[i].Value)
		}
	}
}
