package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	// Setup
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
		pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError) // Reset flags
		initialized = false
		config = EnvConfig{}
	}()

	t.Run("DefaultValues", func(t *testing.T) {
		os.Args = []string{"cmd"}                                           // Reset args to empty
		pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError) // Reset flags
		initialized = false

		cfg := NewConfig()
		assert.Equal(t, "localhost:8080", cfg.ServerAddress)
		assert.Equal(t, "", cfg.FileStoragePath)
		assert.False(t, cfg.EnableHTTPS)
	})

	t.Run("EnvironmentVariables", func(t *testing.T) {
		os.Args = []string{"cmd"}                                           // Reset args to empty
		pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError) // Reset flags
		initialized = false

		os.Setenv("SERVER_ADDRESS", "env:8080")
		os.Setenv("DATABASE_DSN", "env_dsn")
		defer func() {
			os.Unsetenv("SERVER_ADDRESS")
			os.Unsetenv("DATABASE_DSN")
		}()

		cfg := NewConfig()
		assert.Equal(t, "env:8080", cfg.ServerAddress)
		assert.Equal(t, "env_dsn", cfg.DatabaseDSN)
	})

	t.Run("CachedConfig", func(t *testing.T) {
		os.Args = []string{"cmd"}                                           // Reset args to empty
		pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError) // Reset flags
		initialized = false

		cfg1 := NewConfig()
		cfg2 := NewConfig()
		assert.Equal(t, cfg1, cfg2)
	})
}

func TestGetConfigPath(t *testing.T) {
	// Setup
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
		pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError) // Reset flags
	}()

	t.Run("NotSpecified", func(t *testing.T) {
		pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError) // Reset flags
		path := getConfigPath()
		assert.Equal(t, "", path)
	})

	t.Run("FromEnvironment", func(t *testing.T) {
		pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError) // Reset flags
		os.Setenv("CONFIG", "/env/path")
		defer os.Unsetenv("CONFIG")

		path := getConfigPath()
		assert.Equal(t, "/env/path", path)
	})

	t.Run("FromFlag", func(t *testing.T) {
		pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError) // Reset flags
		os.Args = []string{"cmd", "-c", "/test/path"}

		path := getConfigPath()
		assert.Equal(t, "/test/path", path)
	})

}
