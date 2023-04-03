// Package integration contains integration tests and utils.
package integration

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/storage"
)

// SetupDB is a utility function to create an in-memory sqlite database for integration tests.
// It ignores any other database configuration settings.
func SetupDB() error {
	config := configuration.New("../../../example", "../../../../example")
	config.Database.Type = "sqlite"
	config.Database.Path = ":memory:"

	err := storage.InitDB(config)
	if err != nil {
		return fmt.Errorf("could not setup database: %w", err)
	}

	return nil
}

// Migrate is a utility function to create the schema based on the given models.
func Migrate(models ...interface{}) error {
	err := storage.Migrate(models...)

	return errors.Wrap(err, "could not migrate all models")
}
