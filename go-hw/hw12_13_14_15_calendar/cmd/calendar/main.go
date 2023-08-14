package main

import (
	"flag"
	"fmt"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/app/calendar"
	"log"

	"github.com/evgen1067/hw12_13_14_15_calendar/internal/config"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/local.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatalf("Error when reading the configuration file: %s", err)
	}
	logg, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalf("Error during logger initialization: %s", err)
	}
	fmt.Println(cfg)
	err = calendar.Run(cfg, logg)
	if err != nil {
		logg.Error(err.Error())
	}
}
