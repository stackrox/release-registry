package storage_test

import (
	"testing"

	"github.com/stackrox/release-registry/tests"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	err := tests.SetupDB()
	assert.NoError(t, err)
}
