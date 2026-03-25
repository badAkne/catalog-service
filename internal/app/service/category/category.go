package mcategory

import (
	"context"

	"github.com/badAkne/catalog-service/internal/app/entity"
	"github.com/badAkne/catalog-service/internal/app/repository"
	rservice "github.com/badAkne/catalog-service/internal/app/service"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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

// TODO Удалить
func (s *service) SomeMethod(ctx context.Context, req entity.RequestCategoryCreate) error {
	return s.repoCategory.InsideTx(ctx, func(ctx context.Context) error {

		guid, err := uuid.NewV7()
		if err != nil {
			return err
		}

		category := entity.Category{
			GUID: guid,
			Name: req.Name,
		}

		log.Debug().Msg("about to rreate category in tx")

		category, err = s.repoCategory.Create(ctx, category)
		if err != nil {
			log.Debug().Msgf("failed to create category in tx: %s", err.Error())
			return err
		}

		log.Debug().Interface("category", category).Msg("Created category in tx")

		log.Debug().Msg("about to get category in tx")
		category, err = s.repoCategory.Get(ctx, guid)
		if err != nil {
			log.Debug().Msgf("failed to get category %s", err.Error())
			return err
		}

		log.Debug().Interface("category", category).Msg("got category")

		return nil
	})
}
