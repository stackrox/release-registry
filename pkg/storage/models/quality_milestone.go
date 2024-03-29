package models

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/stackrox/release-registry/pkg/storage"
	"github.com/stackrox/release-registry/pkg/utils/validate"
)

// ValidateExpectedMetadataKeysAreProvided checks that all required metadata keys are provided
// and no unexpected keys were passed.
func ValidateExpectedMetadataKeysAreProvided(
	qmd *QualityMilestoneDefinition,
	metadata []QualityMilestoneMetadata,
) error {
	missingMetadataKeys := make(map[string]int)
	unexpectedMetadataKeys := []string{}

	for _, s := range qmd.ExpectedMetadataKeys {
		missingMetadataKeys[s] = 1
	}

	for _, m := range metadata {
		if _, ok := missingMetadataKeys[m.Key]; ok {
			delete(missingMetadataKeys, m.Key)
		} else {
			unexpectedMetadataKeys = append(unexpectedMetadataKeys, m.Key)
		}
	}

	if len(missingMetadataKeys) > 0 || len(unexpectedMetadataKeys) > 0 {
		missingKeysList := []string{}
		for key := range missingMetadataKeys {
			missingKeysList = append(missingKeysList, key)
		}

		return fmt.Errorf(
			"not all required metadata keys were provided, missing: %v, unexpected: %v",
			missingKeysList,
			unexpectedMetadataKeys,
		)
	}

	return nil
}

// ApproveQualityMilestone approves a given Release for a given QualityMilestone.
func ApproveQualityMilestone(
	validApproverDomain, tag, milestoneName, approver string, metadata []QualityMilestoneMetadata,
) (*QualityMilestone, error) {
	if err := validate.IsValidActorEmail(validApproverDomain, approver); err != nil {
		return nil, errors.Wrap(err, "not a valid approver")
	}

	release, err := GetRelease(tag, false, false)
	if err != nil {
		return nil, errors.Wrap(err, "could not find release or already rejected")
	}

	qmd, err := GetQualityMilestoneDefinition(milestoneName)
	if err != nil {
		return nil, errors.Wrap(err, "could not find quality milestone definition")
	}

	if err = ValidateExpectedMetadataKeysAreProvided(qmd, metadata); err != nil {
		return nil, err
	}

	qualityMilestone := &QualityMilestone{
		Approver:                   approver,
		Release:                    *release,
		Metadata:                   metadata,
		QualityMilestoneDefinition: *qmd,
	}

	tx := joinQualityMilestonesWithReleasesAndQualityMilestoneDefinitions(storage.DB)
	tx = tx.Where("releases.tag = ?", release.Tag).Where("quality_milestone_definitions.name = ?", qmd.Name)

	result := tx.FirstOrCreate(qualityMilestone)
	if result.Error != nil {
		return nil, result.Error
	}

	log.Infow(
		"release approved for quality milestone",
		"approver", qualityMilestone.Approver,
		"tag", qualityMilestone.Release.Tag,
		"milestone", qualityMilestone.QualityMilestoneDefinition.Name,
	)

	return qualityMilestone, nil
}
