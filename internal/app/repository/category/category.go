package pcategory

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
		if errors.Is(errors.Unwrap(err), sql.ErrNoRows) {
			return entity.Category{}, entity.ErrNotFound
		}

		return *category, fmt.Errorf("unable to get category,%w", err)
	}

	return *category, nil
}

func (r *repoPG) IsExistWithName(ctx context.Context, name string) error {
	category := new(entity.Category)
	err := r.NewSelect().Model(category).Where("name = ?", name).Scan(ctx)
	if err != nil {
		return fmt.Errorf("unable to check does category exist with name: %w", err)
	}

	if category.GUID != uuid.Nil {
		return entity.ErrCategoryAlreadyExists
	}

	return nil
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
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.Field('C') == "23505" {
			return entity.Category{}, entity.ErrCategoryAlreadyExists
		}

		return entity.Category{}, fmt.Errorf("unable to update category: %w", err)
	}

	return category, nil
}

func (r *repoPG) Delete(ctx context.Context, guid uuid.UUID) error {
	res, err := r.NewDelete().Model((*entity.Category)(nil)).Where("guid = ?", guid).Exec(ctx)
	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.Field('C') == "23503" {
			return entity.ErrCategoryHasRelation
		}

		return fmt.Errorf("unable to delete category: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return entity.ErrNotFound
	}

	return nil
}
