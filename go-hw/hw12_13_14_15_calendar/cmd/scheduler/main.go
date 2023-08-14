package main

import (
	"flag"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/config"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/logger"
	"log"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/local.json", "Path to json configuration file")
}

func main() {
	flag.Parse()
	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal("Error when reading the configuration file: " + err.Error())
	}

	logg, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalf("Error during logger initialization: %s", err.Error())
	}

	scheduler.Run(cfg, logg)
}
