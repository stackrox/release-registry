package models_test

import (
	"testing"

	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/tests/integration"
	"github.com/stretchr/testify/assert"
)

const qualityMilestoneName = "QM1"

func setupQualityMilestoneDefinitionTest(t *testing.T) {
	t.Helper()

	err := integration.SetupDB()
	assert.NoError(t, err)
	err = integration.Migrate(&models.QualityMilestoneDefinition{})
	assert.NoError(t, err)
}

func createFakeQualityMilestoneDefinition(t *testing.T) models.QualityMilestoneDefinition {
	t.Helper()

	metadataKeys := []string{"Abc", "Def", "Ghi"}
	qmd, err := models.CreateQualityMilestoneDefinition(qualityMilestoneName, metadataKeys)
	assert.NoError(t, err)

	return *qmd
}

func TestCreateQualityMilestoneDefinition(t *testing.T) {
	setupQualityMilestoneDefinitionTest(t)

	metadataKeys := []string{"Abc", "Def", "GhiJkl"}
	qmd, err := models.CreateQualityMilestoneDefinition(qualityMilestoneName, metadataKeys)
	assert.NoError(t, err)

	assert.Equal(t, qmd.Name, qualityMilestoneName)
	assert.ElementsMatch(t, metadataKeys, qmd.ExpectedMetadataKeys)
}

func TestCreateQualityMilestoneDefinitionWithInvalidMetadataKeysReturnsError(t *testing.T) {
	setupQualityMilestoneDefinitionTest(t)

	metadataKeys := []string{"Abc", "invalid"}
	_, err := models.CreateQualityMilestoneDefinition(qualityMilestoneName, metadataKeys)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid is not a valid metadata key")
}

func TestGetQualityMilestoneDefinition(t *testing.T) {
	setupQualityMilestoneDefinitionTest(t)

	metadataKeys := []string{"Abc", "Def", "Ghi"}
	_, err := models.CreateQualityMilestoneDefinition(qualityMilestoneName, metadataKeys)
	assert.NoError(t, err)

	qmd, err := models.GetQualityMilestoneDefinition(qualityMilestoneName)
	assert.NoError(t, err)
	assert.Equal(t, qmd.Name, qualityMilestoneName)
	assert.ElementsMatch(t, metadataKeys, qmd.ExpectedMetadataKeys)
}

func TestGetUnknownQualityMilestoneDefinitionReturnsError(t *testing.T) {
	setupQualityMilestoneDefinitionTest(t)

	_, err := models.GetQualityMilestoneDefinition("unknown name")
	assert.ErrorContains(t, err, "record not found")
}

func TestListQualityMilestoneDefinition(t *testing.T) {
	setupQualityMilestoneDefinitionTest(t)

	expectedDefinitions := map[string][]string{
		"QM1": {"Abc", "Def", "Ghi"},
		"QM2": {"Jkl", "Mno", "Pqr"},
	}

	for name, expectedMetadataKeys := range expectedDefinitions {
		_, err := models.CreateQualityMilestoneDefinition(name, expectedMetadataKeys)
		assert.NoError(t, err)
	}

	actualDefinitions, err := models.ListQualityMilestoneDefinitions()
	assert.NoError(t, err)

	assert.Len(t, actualDefinitions, 2)

	for _, qmd := range actualDefinitions {
		expectedMetadataKeys, ok := expectedDefinitions[qmd.Name]
		assert.True(t, ok)
		assert.ElementsMatch(t, expectedMetadataKeys, qmd.ExpectedMetadataKeys)
	}
}
