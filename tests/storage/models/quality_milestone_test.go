package models_test

import (
	"testing"

	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/tests"
	"github.com/stretchr/testify/assert"
)

func setupQualityMilestoneTest(t *testing.T) {
	t.Helper()

	err := tests.SetupDB()
	assert.NoError(t, err)
	err = tests.Migrate(
		&models.Release{},
		&models.ReleaseMetadata{},
		&models.QualityMilestoneMetadata{},
		&models.QualityMilestoneDefinition{},
		&models.QualityMilestone{},
	)
	assert.NoError(t, err)
}

func TestApproveQualityMilestone(t *testing.T) {
	setupQualityMilestoneTest(t)
	release := createFakeRelease(t)
	qmd := createFakeQualityMilestoneDefinition(t)

	qualityMilestone, err := models.ApproveQualityMilestone(
		*configuration.New(),
		release.Tag,
		qmd.Name,
		"roxbot@redhat.com",
		[]models.QualityMilestoneMetadata{
			{Key: "Abc", Value: "abc"},
			{Key: "Def", Value: "def"},
			{Key: "Ghi", Value: "ghi"},
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, release.Tag, qualityMilestone.Release.Tag)
	assert.Equal(t, qmd.Name, qualityMilestone.QualityMilestoneDefinition.Name)
}

func TestApproveRejectedQualityMilestoneReturnsError(t *testing.T) {
	setupQualityMilestoneTest(t)
	release := createFakeRelease(t)

	_, err := models.RejectRelease(release.Tag, false)
	assert.NoError(t, err)

	qmd := createFakeQualityMilestoneDefinition(t)

	_, err = models.ApproveQualityMilestone(
		*configuration.New(),
		release.Tag,
		qmd.Name,
		"roxbot@redhat.com",
		[]models.QualityMilestoneMetadata{
			{Key: "Abc", Value: "abc"},
			{Key: "Def", Value: "def"},
			{Key: "Ghi", Value: "ghi"},
		},
	)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "could not find release or already rejected")
}
