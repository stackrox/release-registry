// Package tests contains integration tests and utils.
package tests

import (
	"fmt"

	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/storage"
)

// SetupDB is a utility function to create an in-memory sqlite database for integration tests.
// It ignores any other database configuration settings.
func SetupDB() error {
	config := configuration.New()
	config.Database.Type = "sqlite"
	config.Database.Path = ":memory:"

	err := storage.InitDB(config)
	if err != nil {
		return fmt.Errorf("could not setup database: %w", err)
	}

	return nil
}
