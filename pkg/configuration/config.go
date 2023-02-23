// Package configuration provides global configuration through files, environment variables and command flags
package configuration

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

//nolint:gochecknoglobals
var (
	once   sync.Once
	config Config
)

// Config is the super structure to hold the database configuration.
type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
}

// DatabaseConfig holds the configuration to access the database.
type DatabaseConfig struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path"`
}

func setupConfigLocation(additionalPaths ...string) {
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("example")
	viper.AddConfigPath("../../example")
	viper.AddConfigPath("/etc")

	for _, p := range additionalPaths {
		viper.AddConfigPath(p)
	}
}

// LoadConfig reads the configuration from a given path.
func LoadConfig(additionalPaths ...string) *Config {
	once.Do(func() {
		setupConfigLocation(additionalPaths...)

		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("cannot read config: %w", err))
		}

		err = viper.Unmarshal(&config)
		if err != nil {
			panic(fmt.Errorf("cannot unmarshal config: %w", err))
		}
	})

	return &config
}
