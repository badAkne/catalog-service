package rcategory

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/badAkne/catalog-service/internal/app/entity"
	rhandler "github.com/badAkne/catalog-service/internal/app/handler"
	rservice "github.com/badAkne/catalog-service/internal/app/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type handler struct {
	serviceCategory rservice.Category
}

func NewHandler(serviceCategory rservice.Category) rhandler.Category {
	return &handler{serviceCategory: serviceCategory}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	reqCategory := entity.RequestCategoryCreate{}
	err := json.NewDecoder(r.Body).Decode(&reqCategory)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusBadRequest)
	}

	defer r.Body.Close()

	if reqCategory.Name == "" {
		http.Error(w, "Обязательное поле отсутствует", http.StatusBadRequest)
	}
	resCategory, err := h.serviceCategory.Create(r.Context(), reqCategory)
	if err != nil {

		var serviceErr = entity.ErrCategoryAlreadyExists
		if errors.Is(err, serviceErr) {
			http.Error(w, "Категория с таким названием уже существует", http.StatusConflict)
			return
		}

		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resCategory); err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.Parse(vars["category_guid"])
	if err != nil {
		http.Error(w, "Неправильный формат UUID", http.StatusBadRequest)
		return
	}

	category, err := h.serviceCategory.Get(r.Context(), guid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			http.Error(w, "Категория не найдена", http.StatusNotFound)
			return
		}

		http.Error(w, "Внутренняя ошибка сервера", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(category); err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	categories, err := h.serviceCategory.GetList(r.Context())
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.Parse(vars["category_guid"])
	if err != nil {
		http.Error(w, ("Неправильный формат UUID"), http.StatusBadRequest)
		return
	}

	category := entity.ResponseCategoryCreate{}
	if err = json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	category, err = h.serviceCategory.Update(r.Context(), guid, category.Name)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			http.Error(w, "Категория не найдена", http.StatusNotFound)
			return
		} else if errors.Is(err, entity.ErrCategoryAlreadyExists) {
			http.Error(w, "Категория с таким названием уже существует", http.StatusConflict)
			return
		}

		http.Error(w, "Внутренняя ошибка сервера", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(category); err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.Parse(vars["category_guid"])
	if err != nil {
		http.Error(w, "Неверный формат UUID", http.StatusBadRequest)
		return
	}

	err = h.serviceCategory.Delete(r.Context(), guid)
	if err != nil {
		if errors.Is(err, entity.ErrCategoryHasRelation) {
			http.Error(w, "Невозможно удалить категорию: имеются связанные товары", http.StatusConflict)
			return
		} else if errors.Is(err, entity.ErrNotFound) {
			http.Error(w, "Категория не найдена", http.StatusBadRequest)
			return
		}

		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
