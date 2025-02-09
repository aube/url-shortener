package main

import (
	"flag"
	"fmt"
)

var serverAddress string
var linkAddress string

func config() {
	flag.StringVar(&serverAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&linkAddress, "b", "http://localhost:8080", "address and port for generated link")

	flag.Parse()

	fmt.Println("serverAddress: " + serverAddress)
}
