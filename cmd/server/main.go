package main

import (
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/logging"
	"github.com/stackrox/release-registry/pkg/service"
)

//nolint:gochecknoglobals
var (
	config = configuration.LoadConfig()
	log    = logging.CreateProductionLogger()
)

func main() {
	log.Infow("Hello from main", "database-type", config.Database.Type, "database-path", config.Database.Path)
	service.HelloWorld()
}
