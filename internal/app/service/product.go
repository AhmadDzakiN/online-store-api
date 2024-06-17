package service

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"online-store-api/internal/app/cache"
	"online-store-api/internal/app/constants"
	"online-store-api/internal/app/payloads"
	"online-store-api/internal/app/repository"
	"online-store-api/internal/pkg/pagination"
)

type ProductService struct {
	Validate    *validator.Validate
	Cache       cache.ICache
	Config      *viper.Viper
	ProductRepo repository.IProductRepository
}

type IProductService interface {
	ViewByCategoryID(ctx echo.Context, categoryID, pageToken string) (resp []payloads.ViewProductResponse, nextPageToken string, err error)
}

func NewProductService(validate *validator.Validate, config *viper.Viper, productRepo repository.IProductRepository, cache cache.ICache) IProductService {
	return &ProductService{
		Validate:    validate,
		Cache:       cache,
		Config:      config,
		ProductRepo: productRepo,
	}
}

func (s *ProductService) ViewByCategoryID(ctx echo.Context, categoryID, pageToken string) (resp []payloads.ViewProductResponse, nextPageToken string, err error) {
	lastValue := pagination.ParsePageToken(pageToken)
	cacheKey := s.Config.GetString("CACHE_ITEM_KEY_VIEW_PRODUCTS_BY_CATEGORY_ID")
	err = s.Cache.Get(ctx.Request().Context(), fmt.Sprintf("%s:%s:last_value:%d", cacheKey, categoryID, lastValue), &resp)
	if err != nil { // Failed to get cache in redis, so no need to return error and get from database instead
		log.Warn().Err(err).Msgf("Failed to get view products by category id cache for category %s", categoryID)
		err = nil

		products, errGetProducts := s.ProductRepo.GetListByCategoryID(ctx.Request().Context(), categoryID, lastValue)
		if errGetProducts != nil {
			log.Err(errGetProducts).Msgf("Failed to get product list by category id %s", categoryID)
			err = errGetProducts
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
	}

	nextPageToken = pagination.CreatePageToken(resp, constants.LimitDataPerPage)
	cacheTTL := s.Config.GetDuration("CACHE_ITEM_TTL_VIEW_PRODUCTS_BY_CATEGORY_ID")
	err = s.Cache.Write(ctx.Request().Context(), fmt.Sprintf("%s:%s:last_value:%d", cacheKey, categoryID, lastValue), resp, cacheTTL)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to write view products by category id cache for category %s", categoryID)
		return
	}

	return
}
