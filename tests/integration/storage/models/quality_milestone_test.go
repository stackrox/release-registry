package models_test

import (
	"testing"

	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/tests/integration"
	"github.com/stackrox/release-registry/tests/utils"
	"github.com/stretchr/testify/assert"
)

func setupQualityMilestoneTest(t *testing.T) {
	t.Helper()

	err := integration.SetupDB()
	assert.NoError(t, err)
	err = integration.Migrate(
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

	qualityMilestoneMetadata := []models.QualityMilestoneMetadata{
		{Key: "Abc", Value: "abc"},
		{Key: "Def", Value: "def"},
		{Key: "Ghi", Value: "ghi"},
	}
	qualityMilestone, err := models.ApproveQualityMilestone(
		configuration.New().Tenant.EmailDomain,
		release.Tag,
		qmd.Name,
		"roxbot@redhat.com",
		qualityMilestoneMetadata,
	)
	assert.NoError(t, err)
	assert.Equal(t, release.Tag, qualityMilestone.Release.Tag)
	assert.Equal(t, qmd.Name, qualityMilestone.QualityMilestoneDefinition.Name)

	approvedRelease, err := models.GetRelease(qualityMilestone.Release.Tag, true, false)
	assert.NoError(t, err)
	utils.AssertReleasesAreEqual(t, &qualityMilestone.Release, approvedRelease, true, true)
	assert.Equal(t, qmd.Name, approvedRelease.QualityMilestones[0].QualityMilestoneDefinition.Name)
}

func TestApproveUnknownReleaseReturnsError(t *testing.T) {
	setupQualityMilestoneTest(t)

	_, err := models.ApproveQualityMilestone(
		configuration.New().Tenant.EmailDomain,
		"1.1.1", "does not matter", "roxbot@redhat.com",
		[]models.QualityMilestoneMetadata{},
	)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "could not find release or already rejected: record not found")
}

func TestApproveUnknownQualityMilestoneDefinitionReturnsError(t *testing.T) {
	setupQualityMilestoneTest(t)

	release := createFakeRelease(t)

	_, err := models.ApproveQualityMilestone(
		configuration.New().Tenant.EmailDomain,
		release.Tag,
		"unknown QualityMilestoneDefinition name",
		"doesnotmatter@redhat.com",
		[]models.QualityMilestoneMetadata{},
	)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "could not find quality milestone definition: record not found")
}

func TestApprovingRejectedQualityMilestoneReturnsError(t *testing.T) {
	setupQualityMilestoneTest(t)
	release := createFakeRelease(t)

	_, err := models.RejectRelease(release.Tag, false)
	assert.NoError(t, err)

	qmd := createFakeQualityMilestoneDefinition(t)

	_, err = models.ApproveQualityMilestone(
		configuration.New().Tenant.EmailDomain,
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
