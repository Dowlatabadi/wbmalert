package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"
)

type Configuration struct {
	Interval int
	Websites []Website
}

func readConfiguration() Configuration {
	var cFlag = flag.String("c", "configuration.json", "path to configuration file")
	flag.Parse()

	jsonFile, err := os.Open(*cFlag)
	if err != nil {
		panic(err)
	}
	byteValue, _ := io.ReadAll(jsonFile)

	var configuration Configuration
	json.Unmarshal(byteValue, &configuration)

	defer jsonFile.Close()

	return configuration
}