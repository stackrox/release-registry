package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/service"
	"github.com/stackrox/release-registry/pkg/storage"
)

func main() {
	config := configuration.New()

	err := storage.InitDB(config)
	if err != nil {
		panic(err)
	}

	err = storage.MigrateAll()
	if err != nil {
		panic(err)
	}

	errCh, err := service.Run(config)
	if err != nil {
		panic(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		panic(errors.Wrap(err, "server error received"))
	case <-sigCh:
		panic(errors.New("signal caught"))
	}
}
