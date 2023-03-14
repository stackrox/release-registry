// Package configuration provides global configuration through files, environment variables and command flags
package configuration

import (
	"fmt"
	"strings"
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
	Type     string `mapstructure:"type"`
	Path     string `mapstructure:"path"`
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Port     int    `mapstructure:"port"`
}

func setupConfigLocation() {
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("example")
	viper.AddConfigPath("../../example")
	viper.AddConfigPath("/etc")
}

func enableEnvVarOverride() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("RELREG")
}

// New is used to generate a configuration instance to pass around the app
func New() *Config {
	once.Do(func() {
		setupConfigLocation()
		enableEnvVarOverride()

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
