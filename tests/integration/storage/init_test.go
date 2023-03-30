package storage_test

import (
	"testing"

	"github.com/stackrox/release-registry/tests/integration"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	err := integration.SetupDB()
	assert.NoError(t, err)
}
