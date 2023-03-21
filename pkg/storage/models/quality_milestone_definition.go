package models

import (
	"fmt"

	"github.com/stackrox/release-registry/pkg/storage"
	"github.com/stackrox/release-registry/pkg/utils/validate"
)

func isMetadataKeyValid(key string) bool {
	return validate.IsValidString(`^([A-Z][a-z]*)+$`, key)
}

func isQualityMilestoneDefinitionNameValid(name string) bool {
	return validate.IsValidString(`^[a-zA-Z0-9 ]+$`, name)
}

// CreateQualityMilestoneDefinition a new QualityMilestoneDefinition.
func CreateQualityMilestoneDefinition(name string, expectedMetadataKeys []string) (*QualityMilestoneDefinition, error) {
	if !isQualityMilestoneDefinitionNameValid(name) {
		return nil, fmt.Errorf("%s is not a valid QualityMilestoneDefinition name", name)
	}

	for _, key := range expectedMetadataKeys {
		if !isMetadataKeyValid(key) {
			return nil, fmt.Errorf("%s is not a valid metadata key", key)
		}
	}

	qmd := &QualityMilestoneDefinition{
		Name:                 name,
		ExpectedMetadataKeys: expectedMetadataKeys,
	}
	result := storage.DB.Where("name = ?", name).FirstOrCreate(qmd)

	if result.Error != nil {
		return nil, result.Error
	}

	log.Infow("quality milestone definition created", "name", qmd.Name)

	return qmd, nil
}

// GetQualityMilestoneDefinition returns a QualityMilestoneDefinition for the given name.
func GetQualityMilestoneDefinition(name string) (*QualityMilestoneDefinition, error) {
	qmd := &QualityMilestoneDefinition{}
	result := storage.DB.Where("name = ?", name).First(qmd)

	if result.Error != nil {
		return nil, result.Error
	}

	return qmd, nil
}

// ListQualityMilestoneDefinitions returns all known QualityMilestoneDefinitions.
func ListQualityMilestoneDefinitions() ([]QualityMilestoneDefinition, error) {
	qualityMilestoneDefinitions := []QualityMilestoneDefinition{}

	result := storage.DB.Find(&qualityMilestoneDefinitions)
	if result.Error != nil {
		return nil, result.Error
	}

	return qualityMilestoneDefinitions, nil
}
