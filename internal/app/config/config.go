package config

import (
	"flag"
	"log"
	"strings"

	"github.com/aube/url-shortener/internal/logger"
	"github.com/caarlos0/env/v6"
)

type EnvConfig struct {
	BaseURL         string `env:"BASE_URL"`
	ServerAddress   string `env:"SERVER_ADDRESS"`
	ServerHost      string
	ServerPort      string
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	FileStorageDir  string `env:"FILE_STORAGE_DIR"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

var config EnvConfig
var initialized bool = false

func getEnvVariables() EnvConfig {
	var cfg EnvConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func NewConfig() EnvConfig {

	if initialized {
		return config
	}

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

	config = getEnvVariables()

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

	config.ServerHost = strings.Split(config.ServerAddress, ":")[0]
	config.ServerPort = strings.Split(config.ServerAddress, ":")[1]

	logger.Println("serverAddress: " + config.ServerAddress)
	logger.Println("serverHost: " + config.ServerHost)
	logger.Println("serverPort: " + config.ServerPort)

	initialized = true

	return config
}
