package rhealth

import (
	"log"
	"net/http"

	rhandler "github.com/badAkne/catalog-service/internal/app/handler"
)

type handler struct{}

func NewHandler() rhandler.Health {
	return &handler{}
}

func (h *handler) LastCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	log.Printf("%v", err)
}
