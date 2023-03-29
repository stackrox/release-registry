package models

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/storage"
	"github.com/stackrox/release-registry/pkg/utils/version"
	"gorm.io/gorm"
)

func withPreloadedMetadata(db *gorm.DB, preloadMetadata bool) *gorm.DB {
	if preloadMetadata {
		return db.Preload("Metadata")
	}

	return db
}

func withPreloadedQualityMilestones(db *gorm.DB, preloadQualityMilestones bool) *gorm.DB {
	if preloadQualityMilestones {
		return db.Preload("QualityMilestones")
	}

	return db
}

func withIncludedRejectedReleases(db *gorm.DB, includeRejected bool) *gorm.DB {
	if !includeRejected {
		return db.Where("rejected = false")
	}

	return db
}

func joinReleasesWithQualityMilestoneDefinitions(tx *gorm.DB) *gorm.DB {
	return tx.Joins(
		"JOIN quality_milestones ON quality_milestones.release_id = releases.id",
	).Joins(
		//nolint:lll
		"JOIN quality_milestone_definitions ON quality_milestones.quality_milestone_definition_id = quality_milestone_definitions.id",
	)
}

func joinQualityMilestonesWithReleasesAndQualityMilestoneDefinitions(tx *gorm.DB) *gorm.DB {
	return tx.Joins(
		"JOIN releases ON quality_milestones.release_id = releases.id",
	).Joins(
		//nolint:lll
		"JOIN quality_milestone_definitions ON quality_milestones.quality_milestone_definition_id = quality_milestone_definitions.id",
	)
}

func findLatestVersionFromListOfReleases(releases []Release) (string, error) {
	if len(releases) == 0 {
		return "", errors.New("no releases found")
	}

	versions := make([]string, len(releases))
	for i, r := range releases {
		versions[i] = r.Tag
	}

	latest, err := version.LatestVersion(versions)
	if err != nil {
		return "", errors.Wrap(err, "could not identify latest version")
	}

	return latest, nil
}

// MigrateAll runs default migrations for all referenced models.
func MigrateAll() error {
	err := storage.Migrate(
		&QualityMilestoneMetadata{},
		&QualityMilestoneDefinition{},
		&QualityMilestone{},
		&ReleaseMetadata{},
		&Release{},
	)
	if err != nil {
		return errors.Wrap(err, "migration of models failed")
	}

	return nil
}

// ValidateActorHasValidEmail checks if the approver has the expected email domain.
func ValidateActorHasValidEmail(config *configuration.Config, approver string) error {
	expectedEmailDomain := config.Tenant.EmailDomain
	if !strings.HasSuffix(approver, expectedEmailDomain) {
		return fmt.Errorf("approver %s has invalid email domain, expected %s", approver, expectedEmailDomain)
	}

	return nil
}
