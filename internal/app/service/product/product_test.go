package mproduct

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/badAkne/catalog-service/internal/app/entity"
	pproduct "github.com/badAkne/catalog-service/internal/app/repository/product"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type createProductSuite struct {
	suite.Suite

	service     *service
	productRepo *pproduct.MockRepository
	ctx         context.Context
}

func (s *createProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = new(pproduct.MockRepository)
	s.service = NewService(s.productRepo).(*service)
}

func (s *createProductSuite) TearDownTest() {
	s.productRepo.AssertExpectations(s.T())
}
func TestCreateProductSuite(t *testing.T) {
	suite.Run(t, new(createProductSuite))
}

type getProductSuite struct {
	suite.Suite

	service     *service
	productRepo *pproduct.MockRepository
	ctx         context.Context
}

func (s *getProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = new(pproduct.MockRepository)
	s.service = NewService(s.productRepo).(*service)
}

func (s *getProductSuite) TearDownTest() {
	s.productRepo.AssertExpectations(s.T())
}
func TestGetProductSuite(t *testing.T) {
	suite.Run(t, new(getProductSuite))
}

type getProductListSuite struct {
	suite.Suite

	service     *service
	productRepo *pproduct.MockRepository
	ctx         context.Context
}

func (s *getProductListSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = new(pproduct.MockRepository)
	s.service = NewService(s.productRepo).(*service)
}

func (s *getProductListSuite) TearDownTest() {
	s.productRepo.AssertExpectations(s.T())
}

func TestGetProductListSuite(t *testing.T) {
	suite.Run(t, new(getProductListSuite))
}

type updateProductSuite struct {
	suite.Suite

	service     *service
	productRepo *pproduct.MockRepository
	ctx         context.Context
}

func (s *updateProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = new(pproduct.MockRepository)
	s.service = NewService(s.productRepo).(*service)
}

func (s *updateProductSuite) TearDownTest() {
	s.productRepo.AssertExpectations(s.T())
}

func TestUpdateProductSuite(t *testing.T) {
	suite.Run(t, new(updateProductSuite))
}

type deleteProductSuite struct {
	suite.Suite

	service     *service
	productRepo *pproduct.MockRepository
	ctx         context.Context
}

func (s *deleteProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = new(pproduct.MockRepository)
	s.service = NewService(s.productRepo).(*service)
}

func (s *deleteProductSuite) TearDownTest() {
	s.productRepo.AssertExpectations(s.T())
}
func TestDeleteProductSuite(t *testing.T) {
	suite.Run(t, new(deleteProductSuite))
}

func (s *createProductSuite) TestCreateProduct() {
	type want struct {
		err    error
		result entity.ResponseProductCreate
	}

	type args struct {
		req entity.RequestProductCreate
	}

	categoryGUID := uuid.New()

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Price:        100.50,
					Description:  "Test Description",
					CategoryGUID: categoryGUID,
				},
			},
			want: want{
				err: nil,
				result: entity.ResponseProductCreate{
					Name:         "Test Product",
					Price:        100.50,
					CategoryGUID: categoryGUID,
					Description:  "Test Description",
				},
			},
			prepare: func(args args) {
				s.productRepo.On("IsExistWithName", s.ctx, args.req.Name).Return(nil).Once()
				createdProduct := entity.Product{
					ID:           123,
					GUID:         uuid.New(),
					Name:         args.req.Name,
					Description:  args.req.Description,
					CategoryGUID: args.req.CategoryGUID,
					Price:        args.req.Price,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}
				s.productRepo.On("Create", s.ctx, mock.MatchedBy(func(p entity.Product) bool {
					return p.Name == args.req.Name &&
						p.Price == args.req.Price &&
						p.CategoryGUID == args.req.CategoryGUID &&
						p.Description == args.req.Description &&
						p.GUID != uuid.Nil
				})).Return(createdProduct, nil).Once()
				s.productRepo.On("InsideTx", s.ctx, mock.AnythingOfType("func(context.Context) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(context.Context) error)
						fn(s.ctx)
					}).
					Return(nil).Once()
			},
		},
		{
			name: "product already exists",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Price:        50.00,
					Description:  "Test ",
					CategoryGUID: categoryGUID,
				},
			},
			want: want{
				err: entity.ErrProductAlreadyExists,
			},
			prepare: func(args args) {
				s.productRepo.On("IsExistWithName", s.ctx, args.req.Name).Return(entity.ErrProductAlreadyExists).Once()
				s.productRepo.On("InsideTx", s.ctx, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(context.Context) error)
					fn(s.ctx)
				}).Return(entity.ErrProductAlreadyExists).Once()
			},
		},
		{
			name: "category doesn't exist",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "New Product",
					Price:        75.00,
					Description:  "New",
					CategoryGUID: categoryGUID,
				},
			},
			want: want{
				err: entity.ErrNotFound,
			},
			prepare: func(args args) {
				s.productRepo.On("IsExistWithName", s.ctx, args.req.Name).Return(nil).Once()
				s.productRepo.On("Create", s.ctx, mock.MatchedBy(func(p entity.Product) bool {
					return p.Name == args.req.Name &&
						p.CategoryGUID == args.req.CategoryGUID
				})).Return(entity.Product{}, entity.ErrNotFound).Once()
				s.productRepo.On("InsideTx", s.ctx, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(context.Context) error)
					fn(s.ctx)
				}).Return(entity.ErrNotFound).Once()
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			testCase.prepare(testCase.args)

			res, err := s.service.Create(s.ctx, testCase.args.req)

			if testCase.want.err != nil {
				s.Error(err, testCase.name)
				s.Empty(res.GUID)
			} else {
				s.NoError(err, testCase.name)
				s.NotEmpty(res.GUID)
				s.Equal(testCase.want.result.Name, res.Name)
				s.Equal(testCase.want.result.Description, res.Description)
				s.Equal(testCase.want.result.CategoryGUID, res.CategoryGUID)
				s.Equal(testCase.want.result.Price, res.Price)
			}
		})
	}
}

func (s *getProductSuite) TestGetProduct() {
	type want struct {
		err error
		res entity.ResponseProductCreate
	}

	type args struct {
		guid uuid.UUID
	}

	now := time.Now()
	guid := uuid.New()
	categoryGUID := uuid.New()

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{
				guid: guid,
			},
			want: want{
				err: nil,
				res: entity.ResponseProductCreate{
					Name:         "Test Product",
					GUID:         guid,
					Description:  "Test Description",
					Price:        100.50,
					CategoryGUID: categoryGUID,
				},
			},
			prepare: func(args args) {
				expectedProduct := entity.Product{
					ID:           1,
					GUID:         guid,
					Name:         "Test Product",
					Description:  "Test Description",
					CategoryGUID: categoryGUID,
					Price:        100.50,
					CreatedAt:    now,
					UpdatedAt:    now,
				}

				s.productRepo.EXPECT().Get(s.ctx, guid).Return(expectedProduct, nil).Once()
			},
		},
		{
			name: "not found",
			args: args{
				guid: uuid.UUID{},
			},
			want: want{
				err: entity.ErrNotFound,
				res: entity.ResponseProductCreate{},
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().Get(s.ctx, args.guid).Return(entity.Product{}, entity.ErrNotFound).Once()
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			testCase.prepare(testCase.args)

			res, err := s.service.Get(s.ctx, testCase.args.guid)

			s.True(errors.Is(err, testCase.want.err), testCase.name)
			s.Equal(testCase.want.res, res)
		})
	}
}

func (s *getProductListSuite) TestGetListProduct() {
	type want struct {
		err error
		res []entity.ResponseProductCreate
	}

	type args struct {
		req entity.RequestProductGetList
	}

	categoryGUID := uuid.New()
	guid1 := uuid.New()
	guid2 := uuid.New()
	guid3 := uuid.New()

	testCases := []struct {
		name    string
		want    want
		args    args
		prepare func(args args)
	}{
		{
			name: "success",
			want: want{
				err: nil,
				res: []entity.ResponseProductCreate{
					{
						Name:         "Test Product 1",
						GUID:         guid1,
						Description:  "Test Description 2",
						Price:        120,
						CategoryGUID: categoryGUID,
					},
					{
						Name:         "Test Product 2",
						GUID:         guid2,
						Description:  "Test Description 2",
						Price:        120,
						CategoryGUID: categoryGUID,
					},
					{
						Name:         "Test Product 3",
						GUID:         guid3,
						Description:  "Test Description 3",
						Price:        130,
						CategoryGUID: categoryGUID,
					},
				},
			},
			args: args{
				req: entity.RequestProductGetList{
					CategoryGUID: categoryGUID,
					MinPrice:     100,
					MaxPrice:     200,
				},
			},
			prepare: func(args args) {
				expectedProducts := []entity.Product{
					{
						Name:         "Test Product 1",
						GUID:         guid1,
						Description:  "Test Description 2",
						Price:        120,
						CategoryGUID: categoryGUID,
					},
					{
						Name:         "Test Product 2",
						GUID:         guid2,
						Description:  "Test Description 2",
						Price:        120,
						CategoryGUID: categoryGUID,
					},
					{
						Name:         "Test Product 3",
						GUID:         guid3,
						Description:  "Test Description 3",
						Price:        130,
						CategoryGUID: categoryGUID,
					},
				}

				s.productRepo.EXPECT().GetList(s.ctx, args.req.CategoryGUID, args.req.MinPrice, args.req.MaxPrice).Return(expectedProducts, nil).Once()
			},
		},
		{
			name: "repo err",
			args: args{
				req: entity.RequestProductGetList{},
			},
			want: want{
				err: entity.ErrNotFound,
				res: nil,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().GetList(s.ctx, args.req.CategoryGUID, args.req.MinPrice, args.req.MaxPrice).Return(nil, entity.ErrNotFound).Once()
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			testCase.prepare(testCase.args)

			res, err := s.service.GetList(s.ctx, testCase.args.req)
			s.True(errors.Is(err, testCase.want.err), testCase.name)
			s.Equal(testCase.want.res, res)
		})
	}

}

func (s *updateProductSuite) TestUpdateProductSuite() {
	type want struct {
		err error
		res entity.ResponseProductCreate
	}

	type args struct {
		guid uuid.UUID
		req  entity.RequestProductUpdate
	}

	guid := uuid.New()
	categoryGuid := uuid.New()
	now := time.Now()

	testCases := []struct {
		name    string
		want    want
		args    args
		prepare func(args args)
	}{
		{
			name: "success",
			want: want{
				err: nil,
				res: entity.ResponseProductCreate{
					GUID:         guid,
					Name:         "Updated product",
					Price:        200,
					CategoryGUID: categoryGuid,
					Description:  "Updated description",
				},
			},
			args: args{
				guid: guid,
				req: entity.RequestProductUpdate{
					Name:         "Updated product",
					Price:        200,
					Description:  "Updated description",
					CategoryGUID: categoryGuid,
				},
			},

			prepare: func(args args) {
				expectedProduct := entity.Product{
					ID:           1,
					GUID:         guid,
					Name:         "Updated product",
					Description:  "Updated description",
					Price:        200,
					CategoryGUID: categoryGuid,
					CreatedAt:    now,
					UpdatedAt:    now,
				}

				s.productRepo.EXPECT().Update(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
					return p.Name == args.req.Name &&
						p.Price == args.req.Price &&
						p.CategoryGUID == args.req.CategoryGUID &&
						p.Description == args.req.Description
				})).Return(expectedProduct, nil).Once()
			},
		},
		{
			name: "not found",
			want: want{
				err: entity.ErrNotFound,
				res: entity.ResponseProductCreate{},
			},
			args: args{
				guid: guid,
				req: entity.RequestProductUpdate{
					Name:         "Test product",
					Price:        200,
					Description:  "Test description",
					CategoryGUID: categoryGuid,
				},
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().Update(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
					return p.GUID == args.guid &&
						p.Name == args.req.Name &&
						p.Price == args.req.Price &&
						p.CategoryGUID == args.req.CategoryGUID &&
						p.Description == args.req.Description
				})).Return(entity.Product{}, entity.ErrNotFound).Once()
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			testCase.prepare(testCase.args)
			res, err := s.service.Update(s.ctx, testCase.args.guid, testCase.args.req)
			s.True(errors.Is(err, testCase.want.err), testCase.name)

			if testCase.want.err == nil {
				s.Equal(guid, res.GUID)
				s.Equal(testCase.args.req.Name, res.Name)
				s.Equal(testCase.args.req.Price, res.Price)
				s.Equal(testCase.args.req.Description, res.Description)
				s.Equal(testCase.args.req.CategoryGUID, res.CategoryGUID)
			} else {
				s.Empty(res.GUID)
			}
		})
	}
}

func (s *deleteProductSuite) TestUpdateProductSuite() {

	type want struct {
		err error
	}

	type args struct {
		guid uuid.UUID
	}

	testCases := []struct {
		name    string
		want    want
		args    args
		prepare func(args args)
	}{
		{
			name: "success",
			want: want{err: nil},
			args: args{
				guid: uuid.New(),
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().Delete(s.ctx, args.guid).Return(nil).Once()
			},
		},
		{
			name: "err not found",
			want: want{
				err: entity.ErrNotFound,
			},
			args: args{
				guid: uuid.New(),
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().Delete(s.ctx, args.guid).Return(entity.ErrNotFound).Once()
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			testCase.prepare(testCase.args)

			err := s.service.Delete(s.ctx, testCase.args.guid)
			s.True(errors.Is(err, testCase.want.err), testCase.name)
		})
	}
}
