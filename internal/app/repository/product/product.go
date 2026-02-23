package pproduct

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/badAkne/catalog-service/internal/app/entity"
	"github.com/badAkne/catalog-service/internal/app/repository"
	rcpostgres "github.com/badAkne/catalog-service/internal/app/repository/conn/postgres"
	"github.com/google/uuid"
	"github.com/uptrace/bun/driver/pgdriver"
)

type (
	repoPg struct {
		*_DB
	}

	_DB = rcpostgres.Client
)

func NewRepoFromPostgres(_ context.Context, d *rcpostgres.Client) (repository.Product, error) {
	return &repoPg{d}, nil
}

func (r *repoPg) Create(ctx context.Context, product entity.Product) (entity.Product, error) {
	_, err := r.NewInsert().Model(&product).Returning("*").Exec(ctx)
	if err != nil {
		var pgErr pgdriver.Error
		if !errors.As(err, &pgErr) {
			return entity.Product{}, fmt.Errorf("unable to create product: %w", err)
		}

		switch pgErr.Field('C') {
		case "23503":
			return entity.Product{}, entity.ErrNotFound
		case "23505":
			return entity.Product{}, entity.ErrProductAlreadyExists
		default:
			return entity.Product{}, fmt.Errorf("unable to create product: %w", err)
		}
	}

	return product, nil
}

func (r *repoPg) Get(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	var product entity.Product

	err := r.NewSelect().Model(&product).Where("guid = ?", guid).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Product{}, entity.ErrNotFound
		}

		return entity.Product{}, fmt.Errorf("unable to get product%w", err)
	}

	return product, nil
}

func (r *repoPg) IsExistWithName(ctx context.Context, name string) error {
	var product entity.Product
	err := r.NewSelect().Model(&product).Where("name = ?", name).Scan(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("unable to check does product exist with name: %w", err)
	}

	if product.GUID != uuid.Nil {
		return entity.ErrProductAlreadyExists
	}

	return nil
}

func (r *repoPg) GetList(ctx context.Context, categoryGUID uuid.UUID, minPrice, maxPrice float32) ([]entity.Product, error) {
	var products []entity.Product

	query := r.NewSelect().Model(&products).Where("price between ? and ?", minPrice, maxPrice)

	err := query.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get products: %w", err)
	}

	return products, nil
}

func (r *repoPg) Update(ctx context.Context, product entity.Product) (entity.Product, error) {
	res, err := r.NewUpdate().Model(&product).OmitZero().Where("guid=", product.GUID).Exec(ctx)
	if rows, _ := res.RowsAffected(); rows == 0 {
		return entity.Product{}, entity.ErrNotFound
	}

	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Field('C') {
			case "23503":
				return entity.Product{}, entity.ErrNotFound
			case "23505":
				return entity.Product{}, entity.ErrProductAlreadyExists
			}
		}

		return entity.Product{}, fmt.Errorf("unable to update product: %w", err)
	}

	return product, nil
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	res, err := r.NewDelete().Model((*entity.Product)(nil)).Where("guid = ?", guid).Exec(ctx)
	if err != nil {
		return fmt.Errorf("unable to delete product: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return entity.ErrNotFound
	}

	return nil
}
