package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
)

var serverAddress string
var baseURL string

type EnvConfig struct {
	BaseURL       string `env:"BASE_URL"`
	ServerAddress string `env:"SERVER_ADDRESS"`
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

	envCfg := getEnvVariables()

	if envCfg.BaseURL > "" {
		baseURL = envCfg.BaseURL
	} else {
		flag.StringVar(&baseURL, "b", "http://localhost:8080", "address and port for generated link")
	}

	if envCfg.ServerAddress > "" {
		serverAddress = envCfg.ServerAddress
	} else {
		flag.StringVar(&serverAddress, "a", "localhost:8080", "address and port to run server")
	}

	flag.Parse()

	fmt.Println("serverAddress: " + serverAddress)
}
