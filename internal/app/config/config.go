package config

import (
	"log"

	"github.com/badAkne/catalog-service/internal/app/config/section"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Monitor    section.Monitor
	Processor  section.Processor
	Repository section.Repository
}

var Root Config

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("%v\n", err)
	}

	err = envconfig.Process("APP", &Root)
	if err != nil {
		log.Fatal(err)
	}
}
