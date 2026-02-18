package rproduct

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/badAkne/catalog-service/internal/app/entity"
	rhandler "github.com/badAkne/catalog-service/internal/app/handler"
	rservice "github.com/badAkne/catalog-service/internal/app/service"
	"github.com/badAkne/catalog-service/internal/app/util"
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
	req := entity.RequestProductCreate{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusInternalServerError)
		return
	}

	res, err := h.serviceProduct.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, util.ErrProductAlreadyExists) {
			http.Error(w, "Товар с таким названием уже существует", http.StatusConflict)
			return
		} else if errors.Is(err, util.ErrCategoryNotFound) {
			http.Error(w, "Категория не найдена", http.StatusNotFound)
			return
		}

		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.Parse(vars["product_guid"])
	if err != nil {
		http.Error(w, "Неправильный формат uuid", http.StatusBadRequest)
		return
	}

	product, err := h.serviceProduct.Get(r.Context(), guid)
	if err != nil {
		http.Error(w, "Товар не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(&product); err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
}
func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	filters := entity.RequestProductGetList{}
	if err := json.NewDecoder(r.Body).Decode(&filters); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	products, err := h.serviceProduct.GetList(r.Context(), filters)
	if err != nil {
		http.Error(w, "Внутрення ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(&products); err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.Parse(vars["product_guid"])
	if err != nil {
		http.Error(w, "Неправильный формат uuid", http.StatusBadRequest)
		return
	}

	req := new(entity.RequestProductCreate)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, "Неправильный формат данных", http.StatusInternalServerError)
		return
	}

	res, err := h.serviceProduct.Update(r.Context(), guid, *req)
	if err != nil {
		if errors.Is(err, util.ErrProductAlreadyExists) {
			http.Error(w, "Товар с таким названием уже существует", http.StatusConflict)
			return
		} else if errors.Is(err, util.ErrProductNotFound) {
			http.Error(w, "Категория не найдена", http.StatusNotFound)
			return
		} else if errors.Is(err, util.ErrProductNotFound) {
			http.Error(w, "Товар не найден", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
}
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.Parse(vars["product_guid"])
	if err != nil {
		http.Error(w, ("Неправильный формат uuid"), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	rows, err := h.serviceProduct.Delete(r.Context(), guid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.Error(w, "Товар не найден", http.StatusNotFound)
	}

	w.WriteHeader(http.StatusNoContent)
}
