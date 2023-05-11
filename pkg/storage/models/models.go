// Package models contains the ORM layer
package models

import (
	"github.com/stackrox/release-registry/pkg/utils/version"
	"gorm.io/gorm"
)

// Release is the main citizen of the registry.
type Release struct {
	gorm.Model
	Tag               string
	Commit            string
	Creator           string
	Metadata          []ReleaseMetadata
	QualityMilestones []QualityMilestone
	Kind              version.Kind
	Rejected          bool
}

// ReleaseMetadata is a key-value struct that can be attached to Releases.
type ReleaseMetadata struct {
	gorm.Model
	Key       string
	Value     string
	Release   Release
	ReleaseID uint
}

// QualityMilestoneMetadata is a key-value struct that can be attached to QualityMilestones.
type QualityMilestoneMetadata struct {
	gorm.Model
	Key                string
	Value              string
	QualityMilestone   QualityMilestone
	QualityMilestoneID uint
}

// QualityMilestoneDefinition is the template for QualityMilestones within the release process.
type QualityMilestoneDefinition struct {
	gorm.Model
	Name                 string
	ExpectedMetadataKeys []string `gorm:"serializer:json"`
}

// QualityMilestone is the progress marker within the release process.
type QualityMilestone struct {
	gorm.Model
	Approver                     string
	Metadata                     []QualityMilestoneMetadata
	QualityMilestoneDefinition   QualityMilestoneDefinition
	Release                      Release
	QualityMilestoneDefinitionID uint
	ReleaseID                    uint
}
