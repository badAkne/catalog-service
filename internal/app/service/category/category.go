package mcategory

import (
	"context"

	"github.com/badAkne/catalog-service/internal/app/entity"
	"github.com/badAkne/catalog-service/internal/app/repository"
	rservice "github.com/badAkne/catalog-service/internal/app/service"
	"github.com/google/uuid"
)

type (
	service struct {
		repoCategory repository.Category
	}
)

func NewService(repoCategory repository.Category) rservice.Category {
	return &service{repoCategory: repoCategory}
}

// pgdriver нужен для того, чтобы ловить ошибку от bun
func (s *service) Create(ctx context.Context, req entity.RequestCategoryCreate) (entity.ResponseCategoryCreate, error) {

	var newCategory entity.Category

	guid, err := uuid.NewV7()
	if err != nil {
		return entity.ResponseCategoryCreate{}, err
	}

	err = s.repoCategory.InsideTx(ctx, func(ctx context.Context) error {

		if err := s.repoCategory.IsExistWithName(ctx, req.Name); err != nil {
			return err
		}

		newCategory, err = s.repoCategory.Create(ctx, entity.Category{Name: req.Name, GUID: guid})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.ResponseCategoryCreate{}, err
	}

	/*
		if err := s.repoCategory.IsExistWithName(ctx, req.Name); err != nil {
			return entity.ResponseCategoryCreate{}, err
		}

		guid, err := uuid.NewV7()
		if err != nil {
			return entity.ResponseCategoryCreate{}, err
		}

		newCategory, err := s.repoCategory.Create(ctx, entity.Category{Name: req.Name, GUID: guid})
		if err != nil {
			return entity.ResponseCategoryCreate{}, err
		}
	*/

	return entity.ResponseCategoryCreate{
		Name:      newCategory.Name,
		GUID:      newCategory.GUID,
		CreatedAt: newCategory.CreatedAt,
	}, nil
}

func (s *service) Get(ctx context.Context, guid uuid.UUID) (entity.ResponseCategoryCreate, error) {
	category, err := s.repoCategory.Get(ctx, guid)
	if err != nil {
		return entity.ResponseCategoryCreate{}, err
	}

	return entity.ResponseCategoryCreate{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
	}, nil
}

func (s *service) GetList(ctx context.Context) ([]entity.ResponseCategoryCreate, error) {
	categories, err := s.repoCategory.GetList(ctx)
	if err != nil {
		return nil, err
	}
	resCategories := make([]entity.ResponseCategoryCreate, 0, len(categories))

	for _, category := range categories {
		resCategory := entity.ResponseCategoryCreate{
			GUID:      category.GUID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
		}
		resCategories = append(resCategories, resCategory)
	}

	return resCategories, nil
}

func (s *service) Update(ctx context.Context, guid uuid.UUID, name string) (entity.ResponseCategoryCreate, error) {
	category, err := s.repoCategory.Update(ctx, guid, name)

	if err != nil {
		return entity.ResponseCategoryCreate{}, err
	}

	if category.GUID == uuid.Nil {
		return entity.ResponseCategoryCreate{}, entity.ErrNotFound
	}

	return entity.ResponseCategoryCreate{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
	}, nil
}

func (s *service) Delete(ctx context.Context, guid uuid.UUID) error {
	err := s.repoCategory.Delete(ctx, guid)
	if err != nil {
		return err
	}

	return nil
}
