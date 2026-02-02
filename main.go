package main

import (
	"log"

	"github.com/badAkne/catalog-service/internal/app/config"
	rhealth "github.com/badAkne/catalog-service/internal/app/handler/health"
	rprocessor "github.com/badAkne/catalog-service/internal/app/processor/http"
)

func main() {
	config.Load()

	cfg := config.Root

	hHandler := rhealth.NewHandler()

	proc := rprocessor.NewHttp(hHandler, cfg.Processor.WebServer)

	if err := proc.Serve(); err != nil {
		log.Fatal(err)
	}
}
