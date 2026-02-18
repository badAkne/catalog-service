package mproduct

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/badAkne/catalog-service/internal/app/entity"
	"github.com/badAkne/catalog-service/internal/app/repository"
	rservice "github.com/badAkne/catalog-service/internal/app/service"
	"github.com/badAkne/catalog-service/internal/app/util"
	"github.com/google/uuid"
	"github.com/uptrace/bun/driver/pgdriver"
)

type (
	service struct {
		repoProduct repository.Product
	}
)

func NewService(repoProduct repository.Product) rservice.Product {
	return &service{repoProduct: repoProduct}
}

func (s *service) Create(ctx context.Context, req entity.RequestProductCreate) (entity.ResponseProductCreate, error) {
	guid, err := uuid.NewV7()
	if err != nil {
		return entity.ResponseProductCreate{}, fmt.Errorf("unable to create guid for product: %w", err)
	}

	newProduct, err := s.repoProduct.Create(ctx, entity.Product{Name: req.Name, Price: req.Price, CategoryGUID: req.CategoryGUID, Description: req.Description, GUID: guid})

	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.Field('C') == "23505" {
			return entity.ResponseProductCreate{}, util.ErrProductAlreadyExists
		} else if errors.As(err, &pgErr) && pgErr.Field('C') == "23503" {
			return entity.ResponseProductCreate{}, util.ErrCategoryNotFound
		}

		return entity.ResponseProductCreate{}, err
	}

	return entity.ResponseProductCreate{
			Name:         newProduct.Name,
			GUID:         newProduct.GUID,
			Price:        newProduct.Price,
			CategoryGUID: newProduct.GUID,
			Description:  newProduct.Description,
		},
		err
}

func (s *service) Get(ctx context.Context, guid uuid.UUID) (entity.ResponseProductCreate, error) {

	product, err := s.repoProduct.Get(ctx, guid)
	if err != nil {
		if errors.Is(errors.Unwrap(err), sql.ErrNoRows) {
			return entity.ResponseProductCreate{}, util.ErrProductNotFound
		}

		return entity.ResponseProductCreate{}, err
	}

	return entity.ResponseProductCreate{
		GUID:         product.GUID,
		Name:         product.Name,
		Price:        product.Price,
		CategoryGUID: product.CategoryGUID,
		Description:  product.Description,
	}, nil
}

func (s *service) GetList(ctx context.Context, req entity.RequestProductGetList) ([]entity.ResponseProductCreate, error) {

	var defMinPrice float32
	req.MinPrice = max(req.MinPrice, defMinPrice)
	if req.MaxPrice == 0 {
		req.MaxPrice = math.MaxFloat32
	}

	products, err := s.repoProduct.GetList(ctx, req.CategoryGUID, req.MinPrice, req.MaxPrice)
	if err != nil {
		return nil, fmt.Errorf("unable to get products: %w", err)
	}

	resProducts := make([]entity.ResponseProductCreate, 0, len(products))

	for _, product := range products {
		resProduct := entity.ResponseProductCreate{
			GUID:         product.GUID,
			Name:         product.Name,
			Price:        product.Price,
			CategoryGUID: product.CategoryGUID,
			Description:  product.Description,
		}

		resProducts = append(resProducts, resProduct)
	}

	return resProducts, nil
}

func (s *service) Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductCreate) (entity.ResponseProductCreate, error) {

	product, err := s.repoProduct.Update(ctx,
		entity.Product{
			GUID:         guid,
			Name:         req.Name,
			Price:        req.Price,
			CategoryGUID: req.CategoryGUID,
			Description:  req.Description,
		},
	)
	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.Field('C') == "23505" {
			return entity.ResponseProductCreate{}, util.ErrProductAlreadyExists
		}

		return entity.ResponseProductCreate{}, err
	}

	if product.GUID == uuid.Nil {
		return entity.ResponseProductCreate{}, util.ErrProductNotFound
	} else if product.CategoryGUID == uuid.Nil {
		return entity.ResponseProductCreate{}, util.ErrCategoryNotFound
	}

	return entity.ResponseProductCreate{
		Name:         product.Name,
		GUID:         product.GUID,
		Description:  product.Description,
		Price:        product.Price,
		CategoryGUID: product.CategoryGUID,
	}, nil
}

func (s *service) Delete(ctx context.Context, guid uuid.UUID) (int64, error) {
	rows, err := s.repoProduct.Delete(ctx, guid)
	if err != nil {
		return 0, err
	}

	return rows, nil
}
