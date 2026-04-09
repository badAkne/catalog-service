package rproduct

import (
	"errors"
	"net/http"

	"github.com/badAkne/catalog-service/internal/app/entity"
	rhandler "github.com/badAkne/catalog-service/internal/app/handler"
	rservice "github.com/badAkne/catalog-service/internal/app/service"
	"github.com/badAkne/catalog-service/internal/pkg/http/binding"
	"github.com/badAkne/catalog-service/internal/pkg/http/httph"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type handler struct {
	serviceProduct rservice.Product
}

func NewHandler(serviceProduct rservice.Product) rhandler.Product {
	return &handler{serviceProduct: serviceProduct}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req entity.RequestProductCreate
	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
		httph.ErrorApply(r, err)
		return
	}

	res, err := h.serviceProduct.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, entity.ErrProductAlreadyExists) {
			httph.SendError(w, http.StatusConflict, err)
			httph.ErrorApply(r, err)
			return
		} else if errors.Is(err, entity.ErrNotFound) {
			httph.SendError(w, http.StatusNotFound, err)
			httph.ErrorApply(r, err)
			return
		}

		httph.SendError(w, http.StatusInternalServerError, err)
		httph.ErrorApply(r, err)
		return
	}

	httph.SendJSON(w, http.StatusCreated, res)
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.Parse(vars["product_guid"])
	if err != nil {
		httph.SendError(w, http.StatusBadRequest, err)
		httph.ErrorApply(r, err)
		return
	}

	product, err := h.serviceProduct.Get(r.Context(), guid)
	if err != nil {
		httph.SendError(w, http.StatusNotFound, err)
		httph.ErrorApply(r, err)
		return
	}

	httph.SendJSON(w, http.StatusOK, product)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	var filters entity.RequestProductGetList
	if err := binding.ScanAndValidateJSON(r, &filters); err != nil {
		httph.SendError(w, http.StatusBadRequest, err)
		httph.ErrorApply(r, err)
		return
	}

	products, err := h.serviceProduct.GetList(r.Context(), filters)
	if err != nil {
		httph.SendError(w, http.StatusInternalServerError, err)
		httph.ErrorApply(r, err)
		return
	}

	httph.SendJSON(w, http.StatusOK, products)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.Parse(vars["product_guid"])
	if err != nil {
		httph.SendError(w, http.StatusBadRequest, err)
		httph.ErrorApply(r, err)
		return
	}

	var req entity.RequestProductUpdate
	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.SendError(w, http.StatusBadRequest, err)
		httph.ErrorApply(r, err)
		return
	}

	res, err := h.serviceProduct.Update(r.Context(), guid, req)
	switch {
	case errors.Is(err, entity.ErrProductAlreadyExists):
		httph.SendError(w, http.StatusConflict, err)
		httph.ErrorApply(r, err)
		return
	case errors.Is(err, entity.ErrNotFound):
		httph.SendError(w, http.StatusNotFound, err)
		httph.ErrorApply(r, err)
		return
	case err != nil:
		httph.SendError(w, http.StatusInternalServerError, err)
		httph.ErrorApply(r, err)
		return
	}

	httph.SendJSON(w, http.StatusOK, res)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.Parse(vars["product_guid"])
	if err != nil {
		httph.SendError(w, http.StatusBadRequest, err)
		httph.ErrorApply(r, err)
		return
	}

	err = h.serviceProduct.Delete(r.Context(), guid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.SendError(w, http.StatusNotFound, err)
			httph.ErrorApply(r, err)
			return
		}

		httph.SendError(w, http.StatusInternalServerError, err)
		httph.ErrorApply(r, err)
		return
	}

	httph.SendEmpty(w, http.StatusNoContent)
}
