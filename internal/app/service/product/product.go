package mproduct

import (
	"context"
	"time"

	"github.com/badAkne/catalog-service/internal/app/entity"
	"github.com/badAkne/catalog-service/internal/app/repository"
	rservice "github.com/badAkne/catalog-service/internal/app/service"
	"github.com/google/uuid"
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
	var newProduct entity.Product

	guid, _ := uuid.NewV7()

	err := s.repoProduct.InsideTx(ctx, func(ctx context.Context) error {
		err := s.repoProduct.IsExistWithName(ctx, req.Name)
		if err != nil {
			return err
		}

		product := entity.Product{
			GUID:         guid,
			Name:         req.Name,
			Price:        req.Price,
			CategoryGUID: req.CategoryGUID,
			Description:  req.Description,
		}

		newProduct, err = s.repoProduct.Create(ctx, product)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return entity.ResponseProductCreate{}, err
	}

	return entity.ResponseProductCreate{
			Name:         newProduct.Name,
			GUID:         newProduct.GUID,
			Price:        newProduct.Price,
			CategoryGUID: newProduct.CategoryGUID,
			Description:  newProduct.Description,
		},
		err
}

func (s *service) Get(ctx context.Context, guid uuid.UUID) (entity.ResponseProductCreate, error) {
	product, err := s.repoProduct.Get(ctx, guid)
	if err != nil {
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
	products, err := s.repoProduct.GetList(ctx, req.CategoryGUID, req.MinPrice, req.MaxPrice)
	if err != nil {
		return nil, err
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

func (s *service) Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.ResponseProductCreate, error) {
	product := entity.Product{
		GUID:         guid,
		Name:         req.Name,
		Price:        req.Price,
		CategoryGUID: req.CategoryGUID,
		Description:  req.Description,
		UpdatedAt:    time.Now(),
	}

	product, err := s.repoProduct.Update(ctx, product)
	if err != nil {
		return entity.ResponseProductCreate{}, err
	}

	return entity.ResponseProductCreate{
		Name:         product.Name,
		GUID:         product.GUID,
		Description:  product.Description,
		Price:        product.Price,
		CategoryGUID: product.CategoryGUID,
	}, nil
}

func (s *service) Delete(ctx context.Context, guid uuid.UUID) error {
	err := s.repoProduct.Delete(ctx, guid)
	if err != nil {
		return err
	}

	return nil
}
