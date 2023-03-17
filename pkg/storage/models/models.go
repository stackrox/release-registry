// Package models contains the ORM layer
package models

import (
	"gorm.io/gorm"
)

// Release is the main citizen of the registry.
type Release struct {
	gorm.Model
	Tag               string
	Commit            string
	Creator           string
	Metadata          []Metadata `gorm:"foreignKey:ResourceID"`
	QualityMilestones []QualityMilestone
	Rejected          bool
}

// Metadata is a key-value struct that can be attached to Releases or QualityMilestones.
type Metadata struct {
	gorm.Model
	Key        string
	Value      string
	ResourceID uint
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
	Metadata                     []Metadata `gorm:"foreignKey:ResourceID"`
	QualityMilestoneDefinition   QualityMilestoneDefinition
	Release                      Release
	QualityMilestoneDefinitionID uint
	ReleaseID                    uint
}
