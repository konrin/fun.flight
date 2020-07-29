package main

import (
	"flag"
	"io/ioutil"
	"log"

	"fun.flight/pkg/app"
	"fun.flight/pkg/config"
)

func main() {
	var configFilePath string

	flag.StringVar(&configFilePath, "config", "./config.yaml", "")
	flag.Parse()

	configFileData, err := readConfigFile(configFilePath)
	if err != nil {
		log.Printf("Error config file: \"%s\"\n", configFilePath)

		return
	}

	conf, err := config.Load(configFileData)
	if err != nil {
		log.Printf("Error load config: \"%s\"\n", err.Error())

		return
	}

	if err := app.Run(conf); err != nil {
		log.Println(err.Error())
	}
}

func readConfigFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}
