package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/logging"
	"github.com/stackrox/release-registry/pkg/service"
	"github.com/stackrox/release-registry/pkg/storage"
	"github.com/stackrox/release-registry/pkg/storage/models"
)

//nolint:gochecknoglobals
var log = logging.CreateProductionLogger()

func main() {
	config := configuration.New("./example")

	err := storage.InitDB(config)
	if err != nil {
		log.Fatalw("received an error on database init", "error", err)
	}

	err = models.MigrateAll()
	if err != nil {
		log.Fatalw("received an error on migration", "error", err)
	}

	server := service.New(config)
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
