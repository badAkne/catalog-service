package repository

import (
	"context"

	"github.com/badAkne/catalog-service/internal/app/entity"
	"github.com/google/uuid"
)

type (
	Category interface {
		Create(ctx context.Context, category entity.Category) (entity.Category, error)
		Get(ctx context.Context, guid uuid.UUID) (entity.Category, error)
		IsExistWithName(ctx context.Context, name string) error
		GetList(ctx context.Context) ([]entity.Category, error)
		Update(ctx context.Context, guid uuid.UUID, name string) (entity.Category, error)
		Delete(ctx context.Context, guid uuid.UUID) (int64, error)
	}
	Product interface {
		Create(ctx context.Context, product entity.Product) (entity.Product, error)
		Get(ctx context.Context, guid uuid.UUID) (entity.Product, error)
		IsExistWithName(ctx context.Context, name string) error
		GetList(ctx context.Context, categoryGUID uuid.UUID, minPrice, maxPrice float32) ([]entity.Product, error)
		Update(ctx context.Context, product entity.Product) (entity.Product, error)
		Delete(ctx context.Context, guid uuid.UUID) (int64, error)
	}
)
