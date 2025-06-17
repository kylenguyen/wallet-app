package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

var ErrEnvVarsNotSet = errors.New("env vars not set")

type Config struct {
	ServiceName string
	Env         string
	ServicePort int

	DatabaseVar DatabaseVar
}

type DatabaseVar struct {
	Name            string
	Host            string
	User            string
	Password        string
	Port            int
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// Load reads the configuration from environment variables and returns a Config struct.
// It uses Viper to automatically load environment variables.
// It also validates the loaded configuration to ensure all required values are present and valid.
//
// Returns:
//   - Config: The loaded configuration.
//   - error: An error if the configuration could not be loaded or is invalid.
func Load() (Config, error) {
	env, ok := os.LookupEnv("ENV")
	if ok && env == "prod" {
		viper.SetConfigFile("./deployments/prod.env")
	} else {
		viper.SetConfigFile("./deployments/local.env")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	config := Config{
		ServiceName: viper.GetString("SERVICE_NAME"),
		Env:         viper.GetString("ENV"),
		ServicePort: viper.GetInt("SERVICE_PORT"),

		DatabaseVar: DatabaseVar{
			Name:            viper.GetString("DB_NAME"),
			Host:            viper.GetString("DB_HOST"),
			User:            viper.GetString("DB_USER"),
			Password:        viper.GetString("DB_PASSWORD"),
			Port:            viper.GetInt("DB_PORT"),
			MaxOpenConns:    viper.GetInt("MAXOPENCONNS"),
			MaxIdleConns:    viper.GetInt("MAXIDLECONNS"),
			ConnMaxLifetime: viper.GetDuration("CONNMAXLIFETIME"),
		},
	}
	if err := config.validate(); err != nil {
		return config, err
	}

	return config, nil
}

func (config Config) validate() error {
	if config.ServiceName == "" {
		return fmt.Errorf("SERVICE_NAME: %w", ErrEnvVarsNotSet)
	}

	if config.ServicePort <= 0 {
		return fmt.Errorf("SERVICE_PORT: %w", ErrEnvVarsNotSet)
	}

	if config.DatabaseVar.Name == "" {
		return fmt.Errorf("DB_NAME: %w", ErrEnvVarsNotSet)
	}

	if config.DatabaseVar.Host == "" {
		return fmt.Errorf("DB_HOST: %w", ErrEnvVarsNotSet)
	}

	if config.DatabaseVar.User == "" {
		return fmt.Errorf("DB_USER: %w", ErrEnvVarsNotSet)
	}

	if config.DatabaseVar.Password == "" {
		return fmt.Errorf("DB_PASSWORD: %w", ErrEnvVarsNotSet)
	}

	if config.DatabaseVar.Port <= 0 {
		return fmt.Errorf("DB_PORT: %w", ErrEnvVarsNotSet)
	}

	return nil
}
