package rprocessor

import (
	"net/http"

	rhandler "github.com/badAkne/catalog-service/internal/app/handler"
	"github.com/gorilla/mux"
)

func vGenericRegHealthCheck(r *mux.Router, h rhandler.Health) {
	reg(r, http.MethodGet, "/health", http.HandlerFunc(h.LastCheck))
}

func handlerNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
