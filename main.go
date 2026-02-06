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

	pgClient, err := rcpostgres.NewConn(ctx, cfg.Repository.Postgres)
	if err != nil {
		log.Printf("%s", err.Error())
	}

	oldVer, newVer, err := pgClient.Migrate(ctx)
	if err != nil {
		log.Fatal("failed to run migrations")
	}

	if oldVer != newVer {
		log.Printf("old_version:%d\nnew_version:%d\ndatabase migrated\n", oldVer, newVer)
	} else {
		log.Printf("version:%d\ndatabase is up to date ", newVer)
	}

	hHandler := rhealth.NewHandler()

	proc := rprocessor.NewHttp(hHandler, cfg.Processor.WebServer)

	if err := proc.Serve(); err != nil {
		log.Fatal(err)
	}
}
