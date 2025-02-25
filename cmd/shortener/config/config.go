package config

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/caarlos0/env/v6"
)

type EnvConfig struct {
	BaseURL         string `env:"BASE_URL"`
	ServerAddress   string `env:"SERVER_ADDRESS"`
	ServerHost      string
	ServerPort      string
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	FileStorageDir  string `env:"FILE_STORAGE_DIR"`
}

var Config EnvConfig
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
		return Config
	}

	var flagBaseURL string
	var flagServerAddress string
	var flagStoragePath string
	var flagStorageDir string

	flag.StringVar(&flagBaseURL, "b", "http://localhost:8080", "address and port for generated link")
	flag.StringVar(&flagServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagStorageDir, "d", "./_hashes", "hashes dir")
	flag.StringVar(&flagStoragePath, "f", "./_hashes/hashes_list.json", "hashes file")
	flag.Parse()

	Config = getEnvVariables()

	if Config.BaseURL == "" {
		Config.BaseURL = flagBaseURL
	}

	if Config.FileStoragePath == "" {
		Config.FileStoragePath = flagStoragePath
	}

	if Config.FileStorageDir == "" {
		Config.FileStorageDir = flagStorageDir
	}

	if Config.ServerAddress == "" {
		Config.ServerAddress = flagServerAddress
	}

	Config.ServerHost = strings.Split(Config.ServerAddress, ":")[0]
	Config.ServerPort = strings.Split(Config.ServerAddress, ":")[1]

	fmt.Println("serverAddress: " + Config.ServerAddress)
	fmt.Println("serverHost: " + Config.ServerHost)
	fmt.Println("serverPort: " + Config.ServerPort)

	initialized = true

	return Config
}
