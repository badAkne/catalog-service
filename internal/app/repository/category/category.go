package pcategory

import (
	"context"
	"fmt"

	"github.com/badAkne/catalog-service/internal/app/entity"
	"github.com/badAkne/catalog-service/internal/app/repository"
	rcpostgres "github.com/badAkne/catalog-service/internal/app/repository/postgres"
	"github.com/google/uuid"
)

type (
	repoPG struct {
		*_DB
	}

	_DB = rcpostgres.Client
)

func NewRepoFromPostgres(_ context.Context, d *rcpostgres.Client) (repository.Category, error) {
	return &repoPG{d}, nil
}

func (r *repoPG) Create(ctx context.Context, category entity.Category) (entity.Category, error) {
	_, err := r.NewInsert().Model(&category).Returning("guid,name,created_at").Exec(ctx)
	if err != nil {
		return category, fmt.Errorf("unable to create category: %w", err)
	}

	return category, nil
}

func (r *repoPG) Get(ctx context.Context, guid uuid.UUID) (entity.Category, error) {
	category := new(entity.Category)
	err := r.NewSelect().Model(category).Where("guid = ?", guid).Scan(ctx)
	if err != nil {
		return *category, fmt.Errorf("unable to get category,%w", err)
	}

	return *category, nil
}

func (r *repoPG) GetList(ctx context.Context) ([]entity.Category, error) {
	categories := []entity.Category{}

	err := r.NewSelect().Model(&categories).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get list of categories: %w", err)
	}

	return categories, nil
}

// чтобы pg правильно находил категорию, надо делать uuid.Parse
func (r *repoPG) Update(ctx context.Context, guid uuid.UUID, name string) (entity.Category, error) {
	var category entity.Category

	_, err := r.NewUpdate().Model(&category).Set("name = ?, updated_at = NOW()", name).Where("guid = ?", guid).Returning("*").Exec(ctx)
	if err != nil {
		return entity.Category{}, fmt.Errorf("unable to update category: %w", err)
	}

	return category, nil
}

func (r *repoPG) Delete(ctx context.Context, guid uuid.UUID) (int64, error) {
	res, err := r.NewDelete().Model((*entity.Category)(nil)).Where("guid = ?", guid).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("unable to delete category: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("unable to get rows affected: %w", err)
	}

	return rows, nil
}
