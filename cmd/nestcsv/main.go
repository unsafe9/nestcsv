package main

import (
	"flag"
	"github.com/unsafe9/nestcsv"
	"log"
	"strings"
)

func main() {
	var (
		configPath  string
		commandArgs string
	)
	flag.StringVar(&configPath, "c", "nestcsv.yaml", "config file path")
	flag.StringVar(&commandArgs, "a", "", "command arguments")
	flag.Parse()

	args := strings.Split(commandArgs, " ")
	for i := 0; i < len(args); i++ {
		args[i] = strings.TrimSpace(args[i])
	}

	config, err := nestcsv.ParseConfig(configPath, args)
	if err != nil {
		log.Fatalf("parse config: %v", err)
	}

	if err := nestcsv.Generate(config); err != nil {
		log.Fatalf("generate: %v", err)
	}
}
