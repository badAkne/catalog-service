package main

import (
	"context"
	"log"

	"github.com/badAkne/catalog-service/internal/app/config"
	rcategory "github.com/badAkne/catalog-service/internal/app/handler/category"
	rhealth "github.com/badAkne/catalog-service/internal/app/handler/health"
	rproduct "github.com/badAkne/catalog-service/internal/app/handler/product"
	rprocessor "github.com/badAkne/catalog-service/internal/app/processor/http"
	pcategory "github.com/badAkne/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/badAkne/catalog-service/internal/app/repository/postgres"
	pproduct "github.com/badAkne/catalog-service/internal/app/repository/product"
	mcategory "github.com/badAkne/catalog-service/internal/app/service/category"
	mproduct "github.com/badAkne/catalog-service/internal/app/service/product"
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

	categoryRepo, err := pcategory.NewRepoFromPostgres(ctx, pgClient)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	repoProduct, err := pproduct.NewRepoFromPostgres(ctx, pgClient)
	if err != nil {
		log.Fatalf("unable to create product: %s", err.Error())
	}

	categoryService := mcategory.NewService(categoryRepo)
	productService := mproduct.NewService(repoProduct)

	hHandler := rhealth.NewHandler()
	categoryHandler := rcategory.NewHandler(categoryService)
	productHandler := rproduct.NewHandler(productService)

	proc := rprocessor.NewHttp(hHandler, categoryHandler, productHandler, cfg.Processor.WebServer)

	if err := proc.Serve(); err != nil {
		log.Fatal(err)
	}
}
