package utils

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stackrox/infra-auth-lib/auth"
	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/service"
	"github.com/stackrox/release-registry/pkg/storage"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stretchr/testify/assert"
)

const (
	databaseFileMode = fs.FileMode(0644)
	serverStartWait  = 10 * time.Second

	// DefaultUser is the email contained in the test JWT.
	DefaultUser = "roxbot+release-registry-e2e@redhat.com"
)

// DefaultUserJwt returns the default user's JWT used in tests.
func DefaultUserJwt() string {
	value := os.Getenv("RELREG_TEST_TOKEN")
	if value == "" {
		panic("RELREG_TEST_TOKEN is not set")
	}

	return value
}

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

// SetupE2ETest creates a buffered listener, initializes the database and starts the server.
func SetupE2ETest(t *testing.T, databasePath string) {
	t.Helper()

	config := configuration.New("../../../example")

	tmpDatabasePath, err := copyDatabaseFixtureToTmp(databasePath, "e2e-test-")
	assert.NoError(t, err)

	config.Database = configuration.DatabaseConfig{
		Type: "sqlite",
		Path: tmpDatabasePath,
	}

	RemotePort = getRandomIntInRange(minPort, maxPort)
	config.Server.Port = RemotePort
	config.Server.Cert = "../../../example/server.crt"
	config.Server.Key = "../../../example/server.key"

	oidc, err := auth.NewFromConfig("../../../example/oidc.yaml")
	if err != nil {
		panic(err)
	}

	services, err := service.CreateAPIServices(config, *oidc)

	if err != nil {
		panic(err)
	}

	s := service.New(config, *oidc, services...)

	err = storage.InitDB(config)
	if err != nil {
		panic(err)
	}

	err = models.MigrateAll()
	if err != nil {
		panic(err)
	}

	go func() {
		if err := s.Run(); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	time.Sleep(serverStartWait)
}

// GetFixturesPath constructs the absolute path to the fixtures directory.
func GetFixturesPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "could not get working directory")
	}

	return fmt.Sprintf("%s/%s", cwd, "../fixtures"), nil
}
