package service

import (
	"database/sql"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
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
	CartRepo     repository.ICartRepository
	CartItemRepo repository.ICartItemRepository
	ProductRepo  repository.IProductRepository
}

type ICartService interface {
	AddProduct(ctx echo.Context, req payloads.AddProductRequest) (err error)
	DeleteProduct(ctx echo.Context, productID string) (err error)
	View(ctx echo.Context, pageToken string) (resp []payloads.ViewCartResponse, nextPageToken string, err error)
}

func NewCartService(validate *validator.Validate, db *gorm.DB, cartRepo repository.ICartRepository, cartItemRepo repository.ICartItemRepository,
	productRepo repository.IProductRepository) ICartService {
	return &CartService{
		Validate:     validate,
		DB:           db,
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

		// There is already the same product in the cart
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

	return
}

func (s *CartService) DeleteProduct(ctx echo.Context, productID string) (err error) {
	product, err := s.ProductRepo.GetByID(ctx.Request().Context(), productID)
	if err != nil {
		log.Err(err).Msgf("Failed to get product by id %s", productID)
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
	if err != nil {
		log.Err(err).Msgf("Failed to get active cart by customer id %s", customerID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = echo.NewHTTPError(http.StatusNotFound, "cart is not found")
		}
		return
	}

	cartItem, err := s.CartItemRepo.GetActiveByCartIDAndProductID(ctx.Request().Context(), cart.ID, product.ID)
	if err != nil {
		log.Err(err).Msgf("Failed to get cart item by cart id %d", cart.ID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = echo.NewHTTPError(http.StatusNotFound, "cart item is not found")
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

	return
}

func (s *CartService) View(ctx echo.Context, pageToken string) (resp []payloads.ViewCartResponse, nextPageToken string, err error) {
	jwtClaims, err := jwt.GetTokenClaims(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to get jwt token claims")
		err = echo.NewHTTPError(http.StatusUnauthorized, "Failed to get jwt token claims")
		return
	}

	customerID := jwtClaims["customer_id"].(string)
	cart, err := s.CartRepo.GetActiveByCustomerID(ctx.Request().Context(), customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // Don't return error when user do not have any cart
			log.Warn().Msgf("Customer does not have any cart yet")
			err = nil
			return
		}
		log.Err(err).Msgf("Failed to get active cart by customer id %s", customerID)
		return
	}

	lastValue := pagination.ParsePageToken(pageToken)
	items, err := s.CartItemRepo.GetActiveItemsAndProductsByCartID(ctx.Request().Context(), cart.ID, lastValue)
	if err != nil {
		log.Err(err).Msgf("Failed to get active cart items and products cart by cart id %d", cart.ID)
		return
	}

	// No need to return error when user do not have any active cart items
	for _, item := range items {
		resp = append(resp, payloads.ViewCartResponse{
			ProductID:       item.ProductID,
			ProductName:     item.ProductName,
			ProductPrice:    item.ProductPrice,
			ProductQuantity: item.ProductQuantity,
			UpdatedAt:       item.UpdatedAt.Unix(),
		})
	}

	nextPageToken = pagination.CreatePageToken(resp, constants.LimitDataPerPage)

	return
}
