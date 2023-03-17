// Package storage contains database and ORM layers.
package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/logging"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global variable to hold the database access.
//
//nolint:gochecknoglobals
var (
	log = logging.CreateProductionLogger()
	//nolint:varnamelen
	DB *gorm.DB
)

// DBType enumerates the supported databases.
type DBType string

// These const values represent the supported databases.
const (
	SQLite   DBType = "sqlite"
	Postgres DBType = "postgres"
)

var errUnknownDBType = errors.New("unknown db type")

// InitDB is a wrapper around initialization of the selected database type.
func InitDB(config *configuration.Config) error {
	dbType := DBType(config.Database.Type)
	switch dbType {
	case SQLite:
		return initSqlite(config)
	case Postgres:
		return initPostgres(config)
	default:
		return fmt.Errorf("could not init database: %w (%s)", errUnknownDBType, dbType)
	}
}

// Migrate updates the database schema according to the referenced models.
func Migrate(models ...interface{}) error {
	if err := DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	return nil
}

// MigrateAll runs default migrations for all referenced models.
func MigrateAll() error {
	err := Migrate(
		&models.QualityMilestoneDefinition{},
		&models.Metadata{},
	)
	if err != nil {
		return err
	}

	// Apparently these need to run separately to avoid weird errors in Postgres
	err = Migrate(
		&models.QualityMilestone{},
		&models.Release{},
	)
	if err != nil {
		return err
	}

	return nil
}

func ensurePathExists(path string) error {
	directory := filepath.Dir(path)
	if _, err := os.Stat(directory); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(directory, os.ModePerm); err != nil {
			return fmt.Errorf("could not ensure path to database file: %w", err)
		}
	}

	return nil
}

func initSqlite(config *configuration.Config) error {
	var err error

	databasePath := config.Database.Path

	if err = ensurePathExists(databasePath); err != nil {
		return err
	}

	DB, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return fmt.Errorf("could not init sqlite: %w", err)
	}

	log.Infof("opened database at %s", databasePath)

	return nil
}

func initPostgres(config *configuration.Config) error {
	var err error

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password,
		config.Database.Name,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return fmt.Errorf("could not init postgres: %w", err)
	}

	log.Infow(
		"opened postgres database", "database",
		config.Database.Name, "user", config.Database.User,
		"host", config.Database.Host, "port", config.Database.Port,
	)

	return nil
}
