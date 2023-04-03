package models_test

import (
	"testing"

	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/tests/integration"
	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	err := integration.SetupDB()
	assert.NoError(t, err)

	err = models.MigrateAll()
	assert.NoError(t, err)
}
