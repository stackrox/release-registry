// Package service contains only a stub function.
package service

import (
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/logging"
)

//nolint:gochecknoglobals
var log = logging.CreateProductionLogger()

// HelloWorld is a stub function to test log and config initialization from another package.
func HelloWorld(config *configuration.Config) {
	log.Infow("Hello from service", "database-path", config.Database.Path)
}
