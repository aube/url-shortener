package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// EnvConfig holds all configuration parameters for the application from env variables and config file.
type EnvConfig struct {
	BaseURL               string `mapstructure:"base_url" env:"BASE_URL" json:"BASE_URL"`                                              // Base URL for shortened links
	ServerAddress         string `mapstructure:"server_address" env:"SERVER_ADDRESS" json:"SERVER_ADDRESS"`                            // Server address to listen on
	FileStoragePath       string `mapstructure:"file_storage_path" env:"FILE_STORAGE_PATH" json:"FILE_STORAGE_PATH"`                   // Path to file storage
	DatabaseDSN           string `mapstructure:"database_dsn" env:"DATABASE_DSN" json:"DATABASE_DSN"`                                  // Database connection string
	TokenSecretString     string `mapstructure:"token_secret_string" env:"TOKEN_SECRET_STING" json:"TOKEN_SECRET_STING"`               // Secret for JWT tokens (string)
	DefaultRequestTimeout int    `mapstructure:"default_request_timeout" env:"DEFAULT_REQUEST_TIMEOUT" json:"DEFAULT_REQUEST_TIMEOUT"` // Default request timeout in seconds
	PublicCertFile        string `mapstructure:"public_cert_file" env:"PUBLIC_CERT_FILE" json:"PUBLIC_CERT_FILE"`
	PrivateCertFile       string `mapstructure:"private_cert_file" env:"PRIVATE_CERT_FILE" json:"PRIVATE_CERT_FILE"`
	EnableHTTPS           bool   `mapstructure:"enable_https" env:"ENABLE_HTTPS" json:"ENABLE_HTTPS"`
}

var config EnvConfig
var initialized bool = false

// NewConfig initializes and returns the application configuration.
// It loads configuration in the following order of precedence:
// 1. Command-line flags (highest priority)
// 2. Environment variables
// 3. Configuration file
// 4. Default values (lowest priority)
//
// The configuration is loaded only once, subsequent calls return the cached configuration.
func NewConfig() EnvConfig {
	if initialized {
		return config
	}

	var configPath = getConfigPath()

	// Set default values
	viper.SetDefault("server_address", "localhost:8080")
	viper.SetDefault("base_url", "http://localhost:8080")
	viper.SetDefault("token_secret_string", "~_^")
	viper.SetDefault("enable_https", false)
	viper.SetDefault("log_level", "info")
	viper.SetDefault("default_request_timeout", 15)

	if configPath != "" {
		// Step 2: Read from config file (lowest priority)
		viper.SetConfigName("config") // Looks for config.yaml/json/toml
		viper.AddConfigPath(configPath)
		viper.ReadInConfig() // Ignore errors (file is optional)
	}

	viper.AutomaticEnv()

	// Define and parse command-line flags
	pflag.StringP("base_url", "b", "", "Base URL for shortened links")
	pflag.StringP("server_address", "a", "", "Server address to listen on")
	pflag.StringP("database_dsn", "d", "", "Database connection string")
	pflag.StringP("file_storage_path", "f", "", "Path to file storage")
	pflag.StringP("config_path", "c", "", "Path to config.json")
	pflag.BoolP("enable_https", "s", false, "Enable HTTPS")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine) // Flags override everything

	// Unmarshal into struct
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}

	initialized = true

	return config
}

// getConfigPath retrieves the configuration file path from command-line flags or environment variables.
// Returns the path if specified, or an empty string if not configured.
func getConfigPath() string {
	var configPath string

	args := os.Args
	for i := range len(args) {
		arg := args[i]
		fmt.Println("arg", arg)
		if arg[:2] == "-c" {
			switch {
			case len(arg) > 2:
				configPath = arg[3:]
			case len(arg) == 2:
				if i+1 < len(args) {
					configPath = args[i+1]
				}
			}
		}
	}

	if configPath == "" {
		configPath = os.Getenv("CONFIG")
	}

	return configPath
}
