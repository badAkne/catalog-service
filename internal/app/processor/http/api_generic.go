package rprocessor

import (
	"net/http"
	"net/http/pprof"
	"strings"

	rhandler "github.com/badAkne/catalog-service/internal/app/handler"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func vGenericRegHealthCheck(r *mux.Router, h rhandler.Health) {
	reg(r, http.MethodGet, "/health", http.HandlerFunc(h.LastCheck))
}

func handlerNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func vGenericRegMetrics(r *mux.Router) {
	reg(r, http.MethodGet, "/metrics", promhttp.Handler())
}

func GenericRegPprof(r *mux.Router) {
	r.HandleFunc("/debug/pprof", pprofIndex)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	r.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
}

func pprofIndex(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "pprof") {
		r.URL.Path += "/"
		http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
		return
	}
	pprof.Index(w, r)
}
