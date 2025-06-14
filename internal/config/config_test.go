package config_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bitbucket.org/ntuclink/ff-order-history-go/internal/config"
)

func TestConfig_Load(t *testing.T) {
	// Set up environment variables for testing
	t.Setenv("SERVICE_NAME", "test-service")
	t.Setenv("DD_ENV", "test")
	t.Setenv("SERVICE_PORT", "8080")
	t.Setenv("DB_NAME", "testdb")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_USER", "testuser")
	t.Setenv("DB_PASSWORD", "testpass")
	t.Setenv("DB_PORT", "3306")
	t.Setenv("MAXOPENCONNS", "10")
	t.Setenv("MAXIDLECONNS", "5")
	t.Setenv("CONNMAXLIFETIME", "30m")

	// Load the configuration
	cfg, err := config.Load()
	require.NoError(t, err)

	// Assert that the configuration is loaded correctly
	assert.Equal(t, "test-service", cfg.ServiceName)
	assert.Equal(t, "test", cfg.Env)
	assert.Equal(t, 8080, cfg.ServicePort)
	assert.Equal(t, "testdb", cfg.DatabaseVar.Name)
	assert.Equal(t, "localhost", cfg.DatabaseVar.Host)
	assert.Equal(t, "testuser", cfg.DatabaseVar.User)
	assert.Equal(t, "testpass", cfg.DatabaseVar.Password)
	assert.Equal(t, 3306, cfg.DatabaseVar.Port)
	assert.Equal(t, 10, cfg.DatabaseVar.MaxOpenConns)
	assert.Equal(t, 5, cfg.DatabaseVar.MaxIdleConns)
	assert.Equal(t, 30*time.Minute, cfg.DatabaseVar.ConnMaxLifetime)
}

func TestLoad_MissingEnvVars(t *testing.T) {
	testCases := []struct {
		name        string
		env         map[string]string
		expectedErr error
	}{
		{
			name: "Missing SERVICE_NAME",
			env: map[string]string{
				"SERVICE_PORT": "8080",
				"DB_NAME":      "testdb",
				"DB_HOST":      "localhost",
				"DB_USER":      "testuser",
				"DB_PASSWORD":  "testpass",
				"DB_PORT":      "3306",
			},
			expectedErr: fmt.Errorf("SERVICE_NAME: %w", config.ErrEnvVarsNotSet),
		},
		{
			name: "Missing SERVICE_PORT",
			env: map[string]string{
				"SERVICE_NAME": "test-service",
				"DB_NAME":      "testdb",
				"DB_HOST":      "localhost",
				"DB_USER":      "testuser",
				"DB_PASSWORD":  "testpass",
				"DB_PORT":      "3306",
			},
			expectedErr: fmt.Errorf("SERVICE_PORT: %w", config.ErrEnvVarsNotSet),
		},
		{
			name: "Missing DB_NAME",
			env: map[string]string{
				"SERVICE_NAME": "test-service",
				"SERVICE_PORT": "8080",
				"DB_HOST":      "localhost",
				"DB_USER":      "testuser",
				"DB_PASSWORD":  "testpass",
				"DB_PORT":      "3306",
			},
			expectedErr: fmt.Errorf("DB_NAME: %w", config.ErrEnvVarsNotSet),
		},
		{
			name: "Missing DB_HOST",
			env: map[string]string{
				"SERVICE_NAME": "test-service",
				"SERVICE_PORT": "8080",
				"DB_NAME":      "testdb",
				"DB_USER":      "testuser",
				"DB_PASSWORD":  "testpass",
				"DB_PORT":      "3306",
			},
			expectedErr: fmt.Errorf("DB_HOST: %w", config.ErrEnvVarsNotSet),
		},
		{
			name: "Missing DB_USER",
			env: map[string]string{
				"SERVICE_NAME": "test-service",
				"SERVICE_PORT": "8080",
				"DB_NAME":      "testdb",
				"DB_HOST":      "localhost",
				"DB_PASSWORD":  "testpass",
				"DB_PORT":      "3306",
			},
			expectedErr: fmt.Errorf("DB_USER: %w", config.ErrEnvVarsNotSet),
		},
		{
			name: "Missing DB_PASSWORD",
			env: map[string]string{
				"SERVICE_NAME": "test-service",
				"SERVICE_PORT": "8080",
				"DB_NAME":      "testdb",
				"DB_HOST":      "localhost",
				"DB_USER":      "testuser",
				"DB_PORT":      "3306",
			},
			expectedErr: fmt.Errorf("DB_PASSWORD: %w", config.ErrEnvVarsNotSet),
		},
		{
			name: "Missing DB_PORT",
			env: map[string]string{
				"SERVICE_NAME": "test-service",
				"SERVICE_PORT": "8080",
				"DB_NAME":      "testdb",
				"DB_HOST":      "localhost",
				"DB_USER":      "testuser",
				"DB_PASSWORD":  "testpass",
			},
			expectedErr: fmt.Errorf("DB_PORT: %w", config.ErrEnvVarsNotSet),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			_, err := config.Load()
			require.Error(t, err)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
