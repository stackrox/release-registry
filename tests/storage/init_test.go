package storage_test

import (
	"testing"

	"github.com/stackrox/release-registry/pkg/storage"
	"github.com/stackrox/release-registry/tests"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Parallel()

	err := tests.SetupDB()
	assert.NoError(t, err)
}

func TestMigration(t *testing.T) {
	t.Parallel()

	err := tests.SetupDB()
	assert.NoError(t, err)

	err = storage.MigrateAll()
	assert.NoError(t, err)
}
