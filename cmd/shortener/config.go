package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/caarlos0/env/v6"
)

type EnvConfig struct {
	BaseURL       string `env:"BASE_URL"`
	ServerAddress string `env:"SERVER_ADDRESS"`
	ServerPort    string `env:"SERVER_PORT"`
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
	var flagServerAddressPort string

	flag.StringVar(&flagBaseURL, "b", "http://localhost:8080", "address and port for generated link")
	flag.StringVar(&flagServerAddressPort, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	config = getEnvVariables()

	if config.BaseURL == "" {
		config.BaseURL = flagBaseURL
	}

	if config.ServerAddress == "" {
		config.ServerAddress = strings.Split(flagServerAddressPort, ":")[0]
	}

	if config.ServerPort == "" {
		config.ServerPort = strings.Split(flagServerAddressPort, ":")[1]
	}

	fmt.Println("serverAddress: " + config.ServerAddress)
	fmt.Println("serverPort: " + config.ServerPort)

	initialized = true

	return config
}
