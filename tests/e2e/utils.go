// Package e2e contains end-to-end tests and utils.
package e2e

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/service"
	"github.com/stackrox/release-registry/pkg/storage"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"
)

const (
	bufSize          = 1024 * 1024
	databaseFileMode = fs.FileMode(0644)
	minPort          = 30000
	maxPort          = 39999
)

//nolint:gochecknoglobals
var lis *bufconn.Listener

func copyDatabaseFixtureToTmp(databasePath, prefix string) (string, error) {
	input, err := os.ReadFile(databasePath)
	if err != nil {
		return "", errors.Wrap(err, "could not read database")
	}

	f, err := os.CreateTemp("/tmp", prefix)
	if err != nil {
		return "", errors.Wrap(err, "could not create tmp file")
	}

	err = os.WriteFile(f.Name(), input, databaseFileMode)
	if err != nil {
		return "", errors.Wrap(err, "could not write to tmp file")
	}

	return f.Name(), nil
}

func getRandomPort(min, max int) int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(max-min+1) + min
}

// SetupE2ETest creates a buffered listener, initializes the database and starts the server.
func SetupE2ETest(t *testing.T, databasePath string) {
	t.Helper()

	lis = bufconn.Listen(bufSize)
	config := configuration.New()

	tmpDatabasePath, err := copyDatabaseFixtureToTmp(databasePath, "e2e-test-")
	assert.NoError(t, err)

	config.Database = configuration.DatabaseConfig{
		Type: "sqlite",
		Path: tmpDatabasePath,
	}

	config.Server.Port = getRandomPort(minPort, maxPort)

	s := service.New(config)

	err = storage.InitDB(config)
	if err != nil {
		panic(err)
	}

	err = models.MigrateAll()
	if err != nil {
		panic(err)
	}

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

// BufDialer dials the buffered listener.
func BufDialer(context.Context, string) (net.Conn, error) {
	//nolint:wrapcheck
	return lis.Dial()
}

// GetFixturesPath constructs the absolute path to the fixtures directory.
func GetFixturesPath() (string, error) {
	// TODO: this is not working properly from other dirs...
	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "could not get working directory")
	}

	return fmt.Sprintf("%s/%s", cwd, "tests/e2e/fixtures"), nil
}
