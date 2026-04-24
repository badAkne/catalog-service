package rprocessor

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/badAkne/catalog-service/internal/app/config/section"
	rhandler "github.com/badAkne/catalog-service/internal/app/handler"
	"github.com/badAkne/catalog-service/internal/app/processor"
	"github.com/badAkne/catalog-service/internal/app/util"
	"github.com/badAkne/catalog-service/internal/pkg/http/httph"
	"github.com/badAkne/catalog-service/internal/pkg/http/mzerolog"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type httpProc struct {
	server http.Server
	addr   string
}

func NewHttp(hHealth rhandler.Health, hCategory rhandler.Category, hProduct rhandler.Product, middlewares []httph.Middleware, cfg section.ProcessorWebServer) *httpProc {
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

		makeErrorMiddleware(),
	)

	r.Use(middlewaresToGorilla(middlewares)...)

	r.NotFoundHandler = http.HandlerFunc(handlerNotFound)

	rV1 := r.PathPrefix("/v1").Subrouter()
	if hCategory != nil {
		v1RegCategoryHandler(rV1, hCategory)
	}
	if hProduct != nil {
		v1RegProductHandler(rV1, hProduct)
	}
	vGenericRegHealthCheck(r, hHealth)
	vGenericRegMetrics(r)

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

func (p *httpProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	var lc net.ListenConfig
	l, err := lc.Listen(ctx, "tcp", p.addr)
	if err != nil {
		log.Fatal().Err(err).Str("listen_addr", p.addr).Msg("Failed to start listening TCP addr for HTTP servver")
		return
	}

	log.Info().Str("listen_addr", p.addr).Msg("Listening of TCP addr for HTTP server has been started")

	go p.serve(l)

	processor.WatchForShutdown(ctx, wg, util.CloserFunc(l.Close))

	processor.WatchForShutdown(ctx, wg, util.NewCloserContextFunc(p.server.Shutdown, context.Background(), 5*time.Second))
}

func (h *httpProc) serve(l net.Listener) {
	_ = h.server.Serve(l)
}
