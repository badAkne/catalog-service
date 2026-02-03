package main

import (
	"context"
	"log"

	"github.com/badAkne/catalog-service/internal/app/config"
	rhealth "github.com/badAkne/catalog-service/internal/app/handler/health"
	rprocessor "github.com/badAkne/catalog-service/internal/app/processor/http"
	rcpostgres "github.com/badAkne/catalog-service/internal/app/repository/postgres"
)

func main() {
	ctx := context.Background()

	config.Load()

	cfg := config.Root

	_, err := rcpostgres.NewConn(ctx, cfg.Repository.Postgres)
	if err != nil {
		log.Printf("%s", err.Error())
	}

	hHandler := rhealth.NewHandler()

	proc := rprocessor.NewHttp(hHandler, cfg.Processor.WebServer)

	if err := proc.Serve(); err != nil {
		log.Fatal(err)
	}
}
