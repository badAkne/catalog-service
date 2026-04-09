package catalog

import (
	"context"
	"errors"

	pb "github.com/badAkne/catalog-service/gen/grpc/catalog/v1"
	"github.com/badAkne/catalog-service/internal/app/entity"
	rservice "github.com/badAkne/catalog-service/internal/app/service"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	pb.UnimplementedCatalogServiceServer
	productService rservice.Product
}

func NewHandler(productService rservice.Product) *Handler {
	return &Handler{
		productService: productService,
	}
}

func (h *Handler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	guid, err := uuid.Parse(req.GetGuid())
	if err != nil {
		log.Debug().Err(err).Msg("unable to parse uuid")
		return nil, err
	}

	product, err := h.productService.Get(ctx, guid)
	if err != nil {
		log.Debug().Err(err).Msg("Unable to get product")
		return nil, err
	}

	pbProduct := convertProductToProto(product)

	return &pb.GetProductResponse{
		Product: pbProduct,
	}, nil
}

func (h *Handler) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	if req.GetGuid() == nil {
		return nil, errors.New("bad request")
	}

	guids := req.GetGuid()

	products := make([]*pb.Product, 0, len(req.GetGuid()))

	for _, guidStr := range guids {
		guid, err := uuid.Parse(guidStr)
		if err != nil {
			return nil, err
		}

		product, err := h.productService.Get(ctx, guid)
		if err != nil {
			return nil, err
		}

		products = append(products, convertProductToProto(product))
	}

	return &pb.GetProductsResponse{
		Products: products,
	}, nil
}

func (h *Handler) CheckProductExists(ctx context.Context, req *pb.CheckProductExistsRequest) (*pb.CheckProductExistsResponse, error) {
	guid, err := uuid.Parse(req.GetGuid())
	if err != nil {
		log.Debug().Err(err).Msg("unable to parse uuid")
		return nil, err
	}

	product, err := h.productService.Get(ctx, guid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return &pb.CheckProductExistsResponse{
				Exists: false,
				Price:  0,
			}, nil
		}

		return nil, err
	}

	return &pb.CheckProductExistsResponse{
		Exists: true,
		Price:  float64(product.Price),
	}, nil
}

func convertProductToProto(product entity.ResponseProductCreate) *pb.Product {
	return &pb.Product{
		Guid:         product.GUID.String(),
		Name:         product.Name,
		Description:  product.Description,
		Price:        float64(product.Price),
		CategoryGUID: product.CategoryGUID.String(),
	}
}
