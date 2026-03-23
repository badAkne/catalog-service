package rprocessor

import (
	"fmt"
	"log"
	"net/http"

	"github.com/badAkne/catalog-service/internal/app/config/section"
	rhandler "github.com/badAkne/catalog-service/internal/app/handler"
	"github.com/badAkne/catalog-service/internal/app/util"
	"github.com/badAkne/catalog-service/internal/pkg/http/httph"
	"github.com/badAkne/catalog-service/internal/pkg/http/mzerolog"
	"github.com/gorilla/mux"
)

type httpProc struct {
	server http.Server
	addr   string
}

func NewHttp(hHealth rhandler.Health, hCategory rhandler.Category, hProduct rhandler.Product, cfg section.ProcessorWebServer) *httpProc {
	r := mux.NewRouter()

	logMW := mzerolog.NewMiddleware(
		mzerolog.WithSkipper(util.IsFilteredWithHttp),

		mzerolog.WithStringExtractor("user_id", func(r *http.Request) string {
			return r.Header.Get("X-User-ID")
		}),

		mzerolog.WithStringExtractor("session_id", func(r *http.Request) string {
			return r.Header.Get("X-Session-ID")
		}),

		mzerolog.WithStringExtractorOnFail("request_id", func(r *http.Request) string {
			return r.Header.Get("X-Request-ID")
		}),

		mzerolog.WithAnyExtractorOnSuccess("content_length", func(r *http.Request) any {
			if r.ContentLength > 0 {
				return r.ContentLength
			}
			return nil
		}),
	)

	r.Use(
		httph.NewErrorMiddleware(),

		logMW,
	)

	r.NotFoundHandler = http.HandlerFunc(handlerNotFound)

	var rV1 = r.PathPrefix("/v1").Subrouter()
	if hCategory != nil {
		v1RegCategoryHandler(rV1, hCategory)
	}
	if hProduct != nil {
		v1RegProductHandler(rV1, hProduct)
	}
	vGenericRegHealthCheck(r, hHealth)

	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {

		path, _ := route.GetPathTemplate()

		methods, _ := route.GetMethods()

		if path == "" && len(methods) == 0 {
			return nil
		}

		log.Printf("path:%s\nmethods:%s\nRegistered API route", path, methods)

		return nil
	})

	s := httpProc{
		server: http.Server{
			Handler:           r,
			Addr:              fmt.Sprintf(":%d", cfg.ListenPort),
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		},
		addr: fmt.Sprintf(":%d", cfg.ListenPort),
	}

	return &s
}

func (h *httpProc) Serve() error {
	err := h.server.ListenAndServe()
	return err
}
