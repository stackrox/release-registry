package models_test

import (
	"testing"

	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stretchr/testify/assert"
)

func TestValidateActorHasValidEmail(t *testing.T) {
	config := *configuration.New()
	config.Tenant.EmailDomain = "@redhat.com"

	// Invalid email
	err := models.ValidateActorHasValidEmail(config, "roxbot@stackrox.com")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "approver roxbot@stackrox.com has invalid email domain, expected @redhat.com")

	// Valid email
	err = models.ValidateActorHasValidEmail(config, "roxbot@redhat.com")
	assert.NoError(t, err)
}

func TestValidateExpectedMetadataKeysAreProvided(t *testing.T) {
	qmd := &models.QualityMilestoneDefinition{
		Name:                 "QM1",
		ExpectedMetadataKeys: []string{"a", "b", "c"},
	}

	// All keys provided
	allKeysProvided := []models.Metadata{
		{Key: "a", Value: "a"},
		{Key: "b", Value: "b"},
		{Key: "c", Value: "c"},
	}
	err := models.ValidateExpectedMetadataKeysAreProvided(qmd, allKeysProvided)
	assert.NoError(t, err)

	// Key missing
	oneKeyMissing := []models.Metadata{
		{Key: "a", Value: "a"},
		{Key: "b", Value: "b"},
	}
	err = models.ValidateExpectedMetadataKeysAreProvided(qmd, oneKeyMissing)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "not all required metadata keys were provided, missing: [c], unexpected: []")

	// Unexpected key
	unexpectedKey := []models.Metadata{
		{Key: "a", Value: "a"},
		{Key: "b", Value: "b"},
		{Key: "c", Value: "c"},
		{Key: "d", Value: "d"},
	}
	err = models.ValidateExpectedMetadataKeysAreProvided(qmd, unexpectedKey)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "not all required metadata keys were provided, missing: [], unexpected: [d]")

	// Key duplicated
	duplicatedKey := []models.Metadata{
		{Key: "a", Value: "a"},
		{Key: "b", Value: "b"},
		{Key: "c", Value: "c"},
		{Key: "a", Value: "a"},
	}
	err = models.ValidateExpectedMetadataKeysAreProvided(qmd, duplicatedKey)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "not all required metadata keys were provided, missing: [], unexpected: [a]")
}
