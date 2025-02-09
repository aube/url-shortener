package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/caarlos0/env/v6"
)

var serverAddress string
var baseURL string

type EnvConfig struct {
	BaseURL       string `env:"BASE_URL"`
	ServerAddress string `env:"SERVER_ADDRESS"`
	ServerPort    string `env:"SERVER_PORT"`
}

func getEnvVariables() EnvConfig {
	var cfg EnvConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func config() {

	flag.StringVar(&baseURL, "b", "http://localhost:8080", "address and port for generated link")
	flag.StringVar(&serverAddress, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	envCfg := getEnvVariables()

	if envCfg.BaseURL > "" {
		baseURL = envCfg.BaseURL
	}

	if envCfg.ServerAddress > "" {
		serverAddress = envCfg.ServerAddress
	}

	if envCfg.ServerPort > "" {
		serverAddress = strings.Split(serverAddress, ":")[0] + ":" + envCfg.ServerPort
	}

	fmt.Println("serverAddress: " + serverAddress)
}
