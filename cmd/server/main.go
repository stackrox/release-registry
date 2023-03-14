package main

import (
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/logging"
	"github.com/stackrox/release-registry/pkg/service"
	"github.com/stackrox/release-registry/pkg/storage"
)

//nolint:gochecknoglobals
var log = logging.CreateProductionLogger()

func main() {
	config := configuration.New()
	log.Infow("Hello from main", "database-type", config.Database.Type, "database-path", config.Database.Path)
	service.HelloWorld(config)

	err := storage.InitDB(config)
	if err != nil {
		panic(err)
	}

	err = storage.MigrateAll()
	if err != nil {
		panic(err)
	}
}
