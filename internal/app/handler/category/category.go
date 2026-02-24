package rcategory

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
	serviceCategory rservice.Category
}

func NewHandler(serviceCategory rservice.Category) rhandler.Category {
	return &handler{serviceCategory: serviceCategory}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var reqCategory entity.RequestCategoryCreate
	if err := binding.ScanAndValidateJSON(r, &reqCategory); err != nil {
		httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
		return
	}

	defer r.Body.Close()

	resCategory, err := h.serviceCategory.Create(r.Context(), reqCategory)
	if err != nil {
		if errors.Is(err, entity.ErrCategoryAlreadyExists) {
			httph.SendError(w, http.StatusBadRequest, err)
			return
		}

		httph.SendError(w, http.StatusInternalServerError, err)
		return
	}

	httph.SendJSON(w, http.StatusCreated, resCategory)
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.Parse(vars["category_guid"])
	if err != nil {
		httph.SendError(w, http.StatusBadRequest, err)
		return
	}

	category, err := h.serviceCategory.Get(r.Context(), guid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.SendError(w, http.StatusNotFound, entity.ErrNotFound)
			return
		}

		httph.SendError(w, http.StatusInternalServerError, err)
		return
	}

	httph.SendJSON(w, http.StatusOK, category)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	categories, err := h.serviceCategory.GetList(r.Context())
	if err != nil {
		httph.SendError(w, http.StatusInternalServerError, err)
		return
	}

	httph.SendJSON(w, http.StatusOK, categories)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	guid, err := uuid.Parse(vars["category_guid"])
	if err != nil {
		httph.SendError(w, http.StatusBadRequest, err)
		return
	}

	var category entity.RequestCategoryCreate
	if err := binding.ScanAndValidateJSON(r, &category); err != nil {
		httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
		return
	}

	defer r.Body.Close()

	newCategory, err := h.serviceCategory.Update(r.Context(), guid, category.Name)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.SendError(w, http.StatusNotFound, err)
			return
		} else if errors.Is(err, entity.ErrCategoryAlreadyExists) {
			httph.SendError(w, http.StatusConflict, err)
			return
		}

		httph.SendError(w, http.StatusInternalServerError, err)
		return
	}

	httph.SendJSON(w, http.StatusOK, newCategory)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.Parse(vars["category_guid"])
	if err != nil {
		httph.SendError(w, http.StatusBadRequest, err)
		return
	}

	err = h.serviceCategory.Delete(r.Context(), guid)
	if err != nil {
		if errors.Is(err, entity.ErrCategoryHasRelation) {
			httph.SendError(w, http.StatusConflict, err)
			return
		} else if errors.Is(err, entity.ErrNotFound) {
			httph.SendError(w, http.StatusBadRequest, err)
			return
		}

		httph.SendError(w, http.StatusInternalServerError, err)
		return
	}

	httph.SendEmpty(w, http.StatusNoContent)
}
