package rprocessor

import (
	"log"
	"net/http"
	"time"

	"github.com/badAkne/catalog-service/internal/app/config/section"
	rhandler "github.com/badAkne/catalog-service/internal/app/handler"
	"github.com/gorilla/mux"
)

type httpProc struct {
	server http.Server
	addr   string
}

func NewHttp(hHealth rhandler.Health, cfg section.ProcessorWebServer) *httpProc {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(handlerNotFound)

	vGenericRegHealthCheck(r, hHealth)

	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {

		tpl, err := route.GetPathTemplate()

		if err != nil {
			return err
		}

		met, err := route.GetMethods()

		if err != nil {
			return err
		}

		log.Printf("%s\n%s\n", tpl, met)

		return nil
	})

	if err != nil {
		log.Printf("%v", err)
	}

	s := httpProc{
		server: http.Server{
			ReadTimeout:       3 * time.Second,
			WriteTimeout:      6 * time.Second,
			Handler:           r,
			Addr:              cfg.ListenPort,
			IdleTimeout:       15 * time.Second,
			ReadHeaderTimeout: 3 * time.Second,
		},
		addr: cfg.ListenPort,
	}

	return &s
}

func (h *httpProc) Serve() error {
	err := h.server.ListenAndServe()
	return err
}
