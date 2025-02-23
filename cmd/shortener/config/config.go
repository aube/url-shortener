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
	ServerPort      string `env:"SERVER_PORT"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
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
	var flagServerAddressPort string
	var flagFilesDir string

	flag.StringVar(&flagBaseURL, "b", "http://localhost:8080", "address and port for generated link")
	flag.StringVar(&flagServerAddressPort, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagFilesDir, "f", "./_hashes", "hashes dir")
	flag.Parse()

	Config = getEnvVariables()

	if Config.BaseURL == "" {
		Config.BaseURL = flagBaseURL
	}

	if Config.FileStoragePath == "" {
		Config.FileStoragePath = flagFilesDir
	}

	if Config.ServerAddress == "" {
		Config.ServerAddress = strings.Split(flagServerAddressPort, ":")[0]
	}

	if Config.ServerPort == "" {
		Config.ServerPort = strings.Split(flagServerAddressPort, ":")[1]
	}

	fmt.Println("serverAddress: " + Config.ServerAddress)
	fmt.Println("serverPort: " + Config.ServerPort)

	initialized = true

	return Config
}
