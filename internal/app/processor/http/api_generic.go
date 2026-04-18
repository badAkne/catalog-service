package rprocessor

import (
	"net/http"

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
