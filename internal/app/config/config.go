package config

import (
	"flag"
	"log"
	"strings"

	"github.com/caarlos0/env/v6"
)

// EnvConfig holds all configuration parameters for the application.
// It supports both environment variables and command-line flags.
type EnvConfig struct {
	BaseURL               string `env:"BASE_URL"`       // Base URL for shortened links
	ServerAddress         string `env:"SERVER_ADDRESS"` // Server address to listen on
	ServerHost            string // Parsed server host
	ServerPort            string // Parsed server port
	FileStoragePath       string `env:"FILE_STORAGE_PATH"` // Path to file storage
	FileStorageDir        string `env:"FILE_STORAGE_DIR"`  // Directory for file storage
	DatabaseDSN           string `env:"DATABASE_DSN"`      // Database connection string
	TokenSecretString     string `env:"DATABASE_DSN ^_^ "` // Secret for JWT tokens (string)
	TokenSecret           []byte // Secret for JWT tokens (bytes)
	DefaultRequestTimeout int    `env:"DEFAULT_REQUEST_TIMEOUT"` // Default request timeout in seconds
}

var config EnvConfig
var initialized bool = false

// getEnvVariables reads configuration from environment variables.
func getEnvVariables() EnvConfig {
	var cfg EnvConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

// NewConfig initializes and returns the application configuration.
// It reads from environment variables first, then falls back to command-line flags,
// and finally to default values if neither is provided.
// The configuration is cached after first initialization.
func NewConfig() EnvConfig {
	if initialized {
		return config
	}

	// Define command-line flags with defaults
	var flagBaseURL string
	var flagServerAddress string
	var flagStoragePath string
	var flagStorageDir string
	var flagDatabaseDSN string

	flag.StringVar(&flagBaseURL, "b", "http://localhost:8080", "address and port for generated link")
	flag.StringVar(&flagServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagDatabaseDSN, "d", "", "Database connection string")
	flag.StringVar(&flagStorageDir, "dir", "./_hashes", "hashes dir")
	flag.StringVar(&flagStoragePath, "f", "./_hashes/hashes_list.json", "hashes file")
	flag.Parse()

	// Get environment variables
	config = getEnvVariables()

	// Apply fallbacks
	if config.BaseURL == "" {
		config.BaseURL = flagBaseURL
	}
	if config.FileStoragePath == "" {
		config.FileStoragePath = flagStoragePath
	}
	if config.FileStorageDir == "" {
		config.FileStorageDir = flagStorageDir
	}
	if config.ServerAddress == "" {
		config.ServerAddress = flagServerAddress
	}
	if config.DatabaseDSN == "" {
		config.DatabaseDSN = flagDatabaseDSN
	}
	if config.DefaultRequestTimeout < 1 {
		config.DefaultRequestTimeout = 5
	}

	// Parse server address
	parts := strings.Split(config.ServerAddress, ":")
	config.ServerHost = parts[0]
	if len(parts) > 1 {
		config.ServerPort = parts[1]
	}

	// Convert token secret to bytes
	config.TokenSecret = []byte(config.TokenSecretString)

	initialized = true

	return config
}
