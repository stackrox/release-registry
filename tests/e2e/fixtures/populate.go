package main

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/storage"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/tests/e2e/utils"
)

const (
	createQmdErrMessage     = "could not create QualityMilestoneDefinition"
	createReleaseErrMessage = "could not create release"
)

func setupDB(databasePath string) error {
	config := configuration.New("example")

	config.Database = configuration.DatabaseConfig{
		Type: "sqlite",
		Path: databasePath,
	}

	if err := storage.InitDB(config); err != nil {
		return errors.Wrap(err, "received an error on database initialization")
	}

	if err := models.MigrateAll(); err != nil {
		return errors.Wrap(err, "received an error on migration")
	}

	return nil
}

func populateQualityMilestoneDefinitions(databasePath string) error {
	if err := setupDB(databasePath); err != nil {
		return err
	}

	_, err := models.CreateQualityMilestoneDefinition(
		"Nightly has passed",
		[]string{"Image", "BuildURL"},
	)
	if err != nil {
		return errors.Wrap(err, createQmdErrMessage)
	}

	_, err = models.CreateQualityMilestoneDefinition(
		"Canary successful",
		[]string{"DeploymentURL"},
	)
	if err != nil {
		return errors.Wrap(err, createQmdErrMessage)
	}

	return nil
}

//nolint:funlen
func populateReleases(databasePath string) error {
	if err := setupDB(databasePath); err != nil {
		return err
	}

	config := configuration.New()
	config.Tenant.EmailDomain = "redhat.com"

	_, err := models.CreateRelease(
		config.Tenant.EmailDomain,
		"3.73.0",
		"5c321f6a5b62920b02d2f68592fc14b3ceac656e",
		"roxbot@redhat.com",
		[]models.ReleaseMetadata{
			{Key: "Link", Value: "https://github.com/stackrox/stackrox/tree/3.73.0"},
		},
	)
	if err != nil {
		return errors.Wrap(err, createReleaseErrMessage)
	}

	_, err = models.ApproveQualityMilestone(
		config.Tenant.EmailDomain,
		"3.73.0",
		"Nightly has passed",
		"lastname@redhat.com",
		[]models.QualityMilestoneMetadata{
			{Key: "Image", Value: "quay.io/rhacs-eng/main:3.73.0"},
			{Key: "BuildURL", Value: "https://github.com/stackrox/stackrox/actions/runs/4067270806"},
		},
	)
	if err != nil {
		return errors.Wrap(err, "could not approve release")
	}

	_, err = models.CreateRelease(
		config.Tenant.EmailDomain,
		"3.74.x-nightly-20230323",
		"78057dba490233f41b4602f2b2e88775ab7fd4c9",
		"roxbot@redhat.com",
		[]models.ReleaseMetadata{
			{Key: "Link", Value: "https://github.com/stackrox/stackrox/tree/3.74.x-nightly-20230323"},
		},
	)
	if err != nil {
		return errors.Wrap(err, createReleaseErrMessage)
	}

	_, err = models.ApproveQualityMilestone(
		config.Tenant.EmailDomain,
		"3.74.x-nightly-20230323",
		"Nightly has passed",
		"lastname@redhat.com",
		[]models.QualityMilestoneMetadata{
			{Key: "Image", Value: "quay.io/rhacs-eng/main:3.74.x-nightly-20230323-amd64"},
			{Key: "BuildURL", Value: "https://github.com/stackrox/stackrox/actions/runs/4497338772"},
		},
	)
	if err != nil {
		return errors.Wrap(err, "could not approve release")
	}

	return nil
}

func main() {
	basePath, err := utils.GetFixturesPath()
	if err != nil {
		log.Fatal(err)
	}

	databasePath := fmt.Sprintf("%s/%s", basePath, "prepopulated-with-qualitymilestonedefinitions.sqlite")
	if err := populateQualityMilestoneDefinitions(databasePath); err != nil {
		log.Fatal(err)
	}

	databasePath = fmt.Sprintf("%s/%s", basePath, "prepopulated-with-approved-releases.sqlite")
	if err := populateQualityMilestoneDefinitions(databasePath); err != nil {
		log.Fatal(err)
	}

	if err := populateReleases(databasePath); err != nil {
		log.Fatal(err)
	}
}
