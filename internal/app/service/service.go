package rservice

import (
	"context"

	"github.com/badAkne/catalog-service/internal/app/entity"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type (
	Category interface {
		Create(ctx context.Context, req entity.RequestCategoryCreate) (entity.ResponseCategoryCreate, error)
		Get(ctx context.Context, guid uuid.UUID) (entity.ResponseCategoryCreate, error)
		GetList(ctx context.Context) ([]entity.ResponseCategoryCreate, error)
		Update(ctx context.Context, guid uuid.UUID, name string) (entity.ResponseCategoryCreate, error)
		Delete(ctx context.Context, guid uuid.UUID) error
	}
	Product interface {
		Create(ctx context.Context, req entity.RequestProductCreate) (entity.ResponseProductCreate, error)
		Get(ctx context.Context, guid uuid.UUID) (entity.ResponseProductCreate, error)
		GetList(ctx context.Context, req entity.RequestProductGetList) ([]entity.ResponseProductCreate, error)
		Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.ResponseProductCreate, error)
		Delete(ctx context.Context, guid uuid.UUID) error
	}

	Metered interface {
		ProvideMetrics(fact promauto.Factory) []entity.MetricObservation
	}
)
