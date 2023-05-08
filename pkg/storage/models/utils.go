package models

import (
	"github.com/pkg/errors"
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

func withPreloadedQualityMilestoneDefinitions(db *gorm.DB, preloadQualityMilestoneDefinitions bool) *gorm.DB {
	if preloadQualityMilestoneDefinitions {
		return db.Preload("QualityMilestones.QualityMilestoneDefinition")
	}

	return db
}

func withIncludedRejectedReleases(db *gorm.DB, includeRejected bool) *gorm.DB {
	if !includeRejected {
		return db.Where("rejected = false")
	}

	return db
}

func withIgnoredReleaseKinds(db *gorm.DB, ignoredKinds []version.Kind) *gorm.DB {
	for _, kind := range ignoredKinds {
		db = db.Where("kind = ?", kind)
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
