package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"online-store-api/internal/app/constants"
	"online-store-api/internal/app/payloads"
	"online-store-api/internal/app/repository"
	"online-store-api/internal/pkg/pagination"
)

type ProductService struct {
	Validate    *validator.Validate
	ProductRepo repository.IProductRepository
}

type IProductService interface {
	ViewByCategoryID(ctx echo.Context, categoryID, pageToken string) (resp []payloads.ViewProductResponse, nextPageToken string, err error)
}

func NewProductService(validate *validator.Validate, productRepo repository.IProductRepository) IProductService {
	return &ProductService{
		Validate:    validate,
		ProductRepo: productRepo,
	}
}

func (s *ProductService) ViewByCategoryID(ctx echo.Context, categoryID, pageToken string) (resp []payloads.ViewProductResponse, nextPageToken string, err error) {
	lastValue := pagination.ParsePageToken(pageToken)
	products, err := s.ProductRepo.GetListByCategoryID(ctx.Request().Context(), categoryID, lastValue)
	if err != nil {
		log.Err(err).Msgf("Failed to get product list by category id %s", categoryID)
		return
	}

	// No need to return error when the product list size is 0
	if len(products) < 1 {
		return
	}

	for _, product := range products {
		resp = append(resp, payloads.ViewProductResponse{
			ID:        product.ID,
			Name:      product.Name,
			Price:     product.Price,
			UpdatedAt: product.UpdatedAt.Unix(),
		})
	}

	nextPageToken = pagination.CreatePageToken(resp, constants.LimitDataPerPage)

	return
}
