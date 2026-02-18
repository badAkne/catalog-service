package rprocessor

import (
	"net/http"

	rhandler "github.com/badAkne/catalog-service/internal/app/handler"
	"github.com/gorilla/mux"
)

func v1RegCategoryHandler(r1 *mux.Router, h rhandler.Category) {
	reg(r1, http.MethodPost, "/category/create", http.HandlerFunc(h.Create))
	reg(r1, http.MethodGet, "/category/{category_guid}", http.HandlerFunc(h.Get))
	reg(r1, http.MethodPost, "/category/list", http.HandlerFunc(h.GetList))
	reg(r1, http.MethodPut, "/category/{category_guid}", http.HandlerFunc(h.Update))
	reg(r1, http.MethodDelete, "/category/{category_guid}", http.HandlerFunc(h.Delete))
}

func v1RegProductHandler(r1 *mux.Router, h rhandler.Product) {
	reg(r1, http.MethodPost, "/product/create", http.HandlerFunc(h.Create))
	reg(r1, http.MethodGet, "/product/{product_guid}", http.HandlerFunc(h.Get))
	reg(r1, http.MethodPost, "/product/list", http.HandlerFunc(h.GetList))
	reg(r1, http.MethodPut, "/product/{product_guid}", http.HandlerFunc(h.Update))
	reg(r1, http.MethodDelete, "/product/{product_guid}", http.HandlerFunc(h.Delete))
}
