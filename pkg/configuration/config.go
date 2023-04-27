// Package configuration provides global configuration through files, environment variables and command flags
package configuration

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

//nolint:gochecknoglobals
var (
	once   sync.Once
	config Config
)

// Config is the super structure to hold the database configuration.
type Config struct {
	Tenant   TenantConfig   `mapstructure:"tenant"`
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
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

// ServerConfig holds the configuration for the server.
type ServerConfig struct {
	Cert      string `mapstructure:"cert"`
	Key       string `mapstructure:"key"`
	StaticDir string `mapstructure:"staticDir"`
	DocsDir   string `mapstructure:"docsDir"`
	Port      int    `mapstructure:"port"`
}

// TenantConfig holds configuration specific to the tenant.
type TenantConfig struct {
	EmailDomain    string `mapstructure:"emailDomain"`
	Password       string `mapstructure:"password"`
	OidcConfigFile string `mapstructure:"oidcConfigFile"`
}

func setupConfigLocation(additionalConfigDirs ...string) {
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("/config")

	for _, additionalConfigDir := range additionalConfigDirs {
		viper.AddConfigPath(additionalConfigDir)
	}
}

func enableEnvVarOverride() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("RELREG")
}

func setDefaults(input interface{}, parts ...string) error {
	ifv := reflect.ValueOf(input)
	ift := reflect.TypeOf(input)

	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)

		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}

		//nolint:exhaustive
		switch v.Kind() {
		case reflect.Struct:
			if err := setDefaults(v.Interface(), append(parts, tv)...); err != nil {
				return err
			}
		default:
			if err := viper.BindEnv(strings.Join(append(parts, tv), ".")); err != nil {
				return errors.Wrapf(err, "could not bind env var %v, %s", parts, tv)
			}
		}
	}

	return nil
}

// New is used to generate a configuration instance to pass around the app.
func New(additionalConfigDirs ...string) *Config {
	once.Do(func() {
		setupConfigLocation(additionalConfigDirs...)
		enableEnvVarOverride()
		err := setDefaults(config)
		if err != nil {
			panic(fmt.Errorf("cannot set defaults: %w", err))
		}

		err = viper.ReadInConfig()
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
