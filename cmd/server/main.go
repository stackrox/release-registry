package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/stackrox/infra-auth-lib/auth"
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/logging"
	"github.com/stackrox/release-registry/pkg/service"
	"github.com/stackrox/release-registry/pkg/storage"
	"github.com/stackrox/release-registry/pkg/storage/models"
)

//nolint:gochecknoglobals
var log = logging.CreateProductionLogger()

func main() {
	flagConfigDir := flag.String("config-dir", "example", "path to configuration dir")
	flag.Parse()

	config := configuration.New(*flagConfigDir)

	oidc, err := auth.NewFromConfig(config.Tenant.OidcConfigFile)
	if err != nil {
		log.Fatalw("failed to load oidc config file", "path", config.Tenant.OidcConfigFile, "error", err)
	}

	err = storage.InitDB(config)
	if err != nil {
		log.Fatalw("received an error on database init", "error", err)
	}

	err = models.MigrateAll()
	if err != nil {
		log.Fatalw("received an error on migration", "error", err)
	}

	services, err := service.CreateAPIServices(config, *oidc)
	if err != nil {
		log.Fatalw("received an error during creation of svc", "error", err)
	}

	server := service.New(config, *oidc, services...)
	if err = server.Run(); err != nil {
		log.Fatalw("received an error on server start", "error", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-server.ErrCh:
		log.Fatalw("server error received", "error", err)
	case signal := <-sigCh:
		log.Fatalw("signal caught", "signal", signal)
	}
}
