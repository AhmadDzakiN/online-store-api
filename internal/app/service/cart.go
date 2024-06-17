package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
	"online-store-api/internal/app/cache"
	"online-store-api/internal/app/constants"
	"online-store-api/internal/app/model"
	"online-store-api/internal/app/payloads"
	"online-store-api/internal/app/repository"
	"online-store-api/internal/pkg/jwt"
	"online-store-api/internal/pkg/pagination"
	"time"
)

type CartService struct {
	Validate     *validator.Validate
	DB           *gorm.DB
	Cache        cache.ICache
	Config       *viper.Viper
	CartRepo     repository.ICartRepository
	CartItemRepo repository.ICartItemRepository
	ProductRepo  repository.IProductRepository
}

type ICartService interface {
	AddProduct(ctx echo.Context, req payloads.AddProductRequest) (err error)
	DeleteProduct(ctx echo.Context, productID string) (err error)
	View(ctx echo.Context, pageToken string) (resp payloads.ViewCartResponse, nextPageToken string, err error)
}

func NewCartService(validate *validator.Validate, db *gorm.DB, cache cache.ICache, config *viper.Viper, cartRepo repository.ICartRepository, cartItemRepo repository.ICartItemRepository,
	productRepo repository.IProductRepository) ICartService {
	return &CartService{
		Validate:     validate,
		DB:           db,
		Cache:        cache,
		Config:       config,
		CartRepo:     cartRepo,
		CartItemRepo: cartItemRepo,
		ProductRepo:  productRepo,
	}
}

func (s *CartService) AddProduct(ctx echo.Context, req payloads.AddProductRequest) (err error) {
	err = s.Validate.Struct(req)
	if err != nil {
		log.Err(err).Msg("Invalid add product to cart request body")
		err = echo.NewHTTPError(http.StatusBadRequest, "invalid add product to cart request body")
		return
	}

	product, err := s.ProductRepo.GetByID(ctx.Request().Context(), req.ProductID)
	if err != nil {
		log.Err(err).Msgf("Failed to get product by id %s", req.ProductID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = echo.NewHTTPError(http.StatusNotFound, "product is not found")
		}
		return
	}

	jwtClaims, err := jwt.GetTokenClaims(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to get jwt token claims")
		err = echo.NewHTTPError(http.StatusUnauthorized, "Failed to get jwt token claims")
		return
	}

	customerID := jwtClaims["customer_id"].(string)
	cart, err := s.CartRepo.GetActiveByCustomerID(ctx.Request().Context(), customerID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { // If cart not found, then create a new one
		log.Err(err).Msgf("Failed to get active cart by customer id %s", customerID)
		return
	}

	// Customer already have a shopping cart so we just need to create or update cart item based on condition
	if cart.ID > 0 {
		cartItem, errCartItem := s.CartItemRepo.GetActiveByCartIDAndProductID(ctx.Request().Context(), cart.ID, product.ID)
		if errCartItem != nil && !errors.Is(errCartItem, gorm.ErrRecordNotFound) { // if cart item not found, then create a new one
			log.Err(err).Msgf("Failed to get cart item by cart id %d", cart.ID)
			return
		}

		// There is already the same product in the cart item
		if cartItem.ID > 0 {
			cartItem.Quantity += req.Quantity
			cartItem.UpdatedAt = time.Now()
			err = s.CartItemRepo.Update(ctx.Request().Context(), cartItem)
			if err != nil {
				log.Err(err).Msgf("Failed to update cart item for id %d", cartItem.ID)
				return
			}
		} else {
			newCartItem := model.CartItem{
				CartID:    cart.ID,
				ProductID: req.ProductID,
				Quantity:  req.Quantity,
			}

			err = s.CartItemRepo.Create(ctx.Request().Context(), &newCartItem, nil)
			if err != nil {
				log.Err(err).Msgf("Failed to create a new cart item for cart id %d", cart.ID)
				return
			}
		}

		// Invalidate cart cache for this user. Will generate new cart cache in view cart endpoint/logic
		cartCacheKeyPattern := fmt.Sprintf("%s:*", s.Config.GetString("CACHE_ITEM_KEY_VIEW_CART_BY_CUSTOMER_ID"))
		err = s.Cache.DeleteByKeyPattern(ctx.Request().Context(), cartCacheKeyPattern)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to invalidate/delete cart cache by customer id %s", customerID)
			return
		}

		return
	}

	// Customer does not have any shopping cart exist, so create a new one
	trx := s.DB.Begin().WithContext(ctx.Request().Context()) // Start db trx
	defer trx.Rollback()

	newCart := model.Cart{
		CustomerID: customerID,
	}

	err = s.CartRepo.Create(ctx.Request().Context(), &newCart, trx)
	if err != nil {
		log.Err(err).Msgf("Failed to create a new cart for customer id %s", customerID)
		return
	}

	newCartItem := model.CartItem{
		CartID:    newCart.ID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	err = s.CartItemRepo.Create(ctx.Request().Context(), &newCartItem, trx)
	if err != nil {
		log.Err(err).Msgf("Failed to create a new cart item for cart id %d", newCart.ID)
		return
	}

	commitRes := trx.Commit()
	if commitRes.Error != nil {
		log.Err(err).Msg("Failed to commit transaction for add product flow")
		return
	}

	// I think no need to write a new cart cache for a new user who has just created its cart

	return
}

func (s *CartService) DeleteProduct(ctx echo.Context, productID string) (err error) {
	product, err := s.ProductRepo.GetByID(ctx.Request().Context(), productID)
	if err != nil {
		log.Err(err).Msgf("Failed to get product by id %s", productID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = echo.NewHTTPError(http.StatusNotFound, "product is not found")
			return
		}
		return
	}

	jwtClaims, err := jwt.GetTokenClaims(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to get jwt token claims")
		err = echo.NewHTTPError(http.StatusUnauthorized, "Failed to get jwt token claims")
		return
	}

	customerID := jwtClaims["customer_id"].(string)
	cart, err := s.CartRepo.GetActiveByCustomerID(ctx.Request().Context(), customerID)
	if err != nil {
		log.Err(err).Msgf("Failed to get active cart by customer id %s", customerID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = echo.NewHTTPError(http.StatusNotFound, "cart is not found")
			return
		}
		return
	}

	cartItem, err := s.CartItemRepo.GetActiveByCartIDAndProductID(ctx.Request().Context(), cart.ID, product.ID)
	if err != nil {
		log.Err(err).Msgf("Failed to get cart item by cart id %d", cart.ID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = echo.NewHTTPError(http.StatusNotFound, "cart item is not found")
			return
		}
		return
	}

	// No need to delete the cart too if the cart does not have any cart item
	now := time.Now()
	cartItem.DeletedAt = sql.NullTime{Time: now, Valid: true}
	cartItem.Quantity = 0
	cartItem.UpdatedAt = now
	err = s.CartItemRepo.Update(ctx.Request().Context(), cartItem)
	if err != nil {
		log.Err(err).Msgf("Failed to delete cart item for id %d", cartItem.ID)
		return
	}

	// Invalidate cart cache for this user. Will generate new cart cache in view cart endpoint/logic
	cartCacheKeyPattern := fmt.Sprintf("%s:*", s.Config.GetString("CACHE_ITEM_KEY_VIEW_CART_BY_CUSTOMER_ID"))
	err = s.Cache.DeleteByKeyPattern(ctx.Request().Context(), cartCacheKeyPattern)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to invalidate/delete cart cache by customer id %s", customerID)
		return
	}

	return
}

func (s *CartService) View(ctx echo.Context, pageToken string) (resp payloads.ViewCartResponse, nextPageToken string, err error) {
	jwtClaims, err := jwt.GetTokenClaims(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to get jwt token claims")
		err = echo.NewHTTPError(http.StatusUnauthorized, "Failed to get jwt token claims")
		return
	}

	var cartItemListResp []payloads.CartItemResponse
	customerID := jwtClaims["customer_id"].(string)
	lastValue := pagination.ParsePageToken(pageToken)
	cacheKey := s.Config.GetString("CACHE_ITEM_KEY_VIEW_CART_BY_CUSTOMER_ID")
	err = s.Cache.Get(ctx.Request().Context(), fmt.Sprintf("%s:%s:last_value:%d", cacheKey, customerID, lastValue), &resp)
	if err != nil { // Failed to get cache in redis, so no need to return error and get from database instead
		log.Warn().Err(err).Msgf("Failed to get view carts by customer id %s", customerID)
		err = nil

		cart, errGetCart := s.CartRepo.GetActiveByCustomerID(ctx.Request().Context(), customerID)
		if errGetCart != nil {
			if errors.Is(errGetCart, gorm.ErrRecordNotFound) { // Don't return error when user do not have any cart
				log.Warn().Msgf("Customer does not have any cart yet")
				err = nil
				return
			}
			log.Err(errGetCart).Msgf("Failed to get active cart by customer id %s", customerID)
			err = errGetCart
			return
		}

		items, errGetCartItem := s.CartItemRepo.GetActiveItemsAndProductsByCartID(ctx.Request().Context(), cart.ID, lastValue)
		if errGetCartItem != nil {
			log.Err(errGetCartItem).Msgf("Failed to get active cart items and products cart by cart id %d", cart.ID)
			err = errGetCartItem
			return
		}

		// No need to return error when user do not have any active cart items
		if len(items) < 1 {
			return
		}

		for _, item := range items {
			cartItemListResp = append(cartItemListResp, payloads.CartItemResponse{
				CartItemID:      item.CartItemID,
				ProductID:       item.ProductID,
				ProductName:     item.ProductName,
				ProductPrice:    item.ProductPrice,
				ProductQuantity: item.ProductQuantity,
				UpdatedAt:       item.UpdatedAt.Unix(),
			})
		}

		resp.CartID = items[0].CartID
		resp.Items = cartItemListResp
	}

	nextPageToken = pagination.CreatePageToken(cartItemListResp, constants.LimitDataPerPage)
	cacheTTL := s.Config.GetDuration("CACHE_ITEM_TTL_VIEW_CART_BY_CUSTOMER_ID")
	err = s.Cache.Write(ctx.Request().Context(), fmt.Sprintf("%s:%s:last_value:%d", cacheKey, customerID, lastValue), resp, cacheTTL)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to write view cart by customer id cache for customer %s", customerID)
		return
	}

	return
}
