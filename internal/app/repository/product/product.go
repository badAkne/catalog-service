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
	_, err := r.NewInsert().Model(&product).Returning("guid,name,price,category_guid,description").Exec(ctx)
	if err != nil {
		var pgErr pgdriver.Error
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Product{}, entity.ErrNotFound
		} else if errors.As(err, &pgErr) && pgErr.Field('C') == "23503" {
			return entity.Product{}, entity.ErrNotFound
		}

		return entity.Product{}, fmt.Errorf("unable to create product: %w", err)
	}

	return product, nil
}

func (r *repoPg) Get(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	product := new(entity.Product)

	err := r.NewSelect().Model(product).Where("guid = ?", guid).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Product{}, entity.ErrNotFound
		}

		return entity.Product{}, fmt.Errorf("unable to get product%w", err)
	}

	return *product, nil
}

func (r *repoPg) IsExistWithName(ctx context.Context, name string) error {
	product := new(entity.Product)
	err := r.NewSelect().Model(product).Where("name = ?", name).Scan(ctx)
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

	if product.CategoryGUID != uuid.Nil {
		_, err := r.NewUpdate().Model(&product).Set("name = ?,price = ?,category_guid=?,description=?,updated_at = NOW()", product.Name, product.Price, product.CategoryGUID, product.Description).Where("guid = ?", product.GUID).Returning("*").Exec(ctx)
		if err != nil {
			return entity.Product{}, fmt.Errorf("unable to update product: %w", err)
		}
	} else {
		_, err := r.NewUpdate().Model(&product).Set("name = ?,price = ?,description=?,updated_at = NOW()", product.Name, product.Price, product.Description).Where("guid = ?", product.GUID).Returning("*").Exec(ctx)
		if err != nil {
			return entity.Product{}, fmt.Errorf("unable to update product: %w", err)
		}
	}

	return product, nil
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) (int64, error) {
	res, err := r.NewDelete().Model((*entity.Product)(nil)).Where("guid = ?", guid).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("unable to delete product: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}
