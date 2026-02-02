package main

import (
	"log"

	"github.com/badAkne/catalog-service/internal/app/config"
)

func main() {
	config.Load()

	cfg := config.Root

	log.Printf("Server will start on port: %d", cfg.Processor.WebServer.ListenPort)
	log.Printf("Database: %s@%s/%s %d %d",
		cfg.Repository.Postgres.Username,
		cfg.Repository.Postgres.Address,
		cfg.Repository.Postgres.Name,
		cfg.Repository.Postgres.ReadTimeout,
		cfg.Repository.Postgres.WriteTimeout)
	log.Printf("Environment: %s, LogLevel: %s",
		cfg.Monitor.Environment,
		cfg.Monitor.LogLevel)
}
