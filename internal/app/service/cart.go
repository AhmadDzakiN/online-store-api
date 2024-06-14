package service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"online-store-api/internal/app/model"
	"online-store-api/internal/app/payloads"
	"online-store-api/internal/app/repository"
	"online-store-api/internal/pkg/jwt"
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
	cart, err := s.CartRepo.GetByCustomerID(ctx.Request().Context(), customerID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { // If cart not found, then create a new one
		log.Err(err).Msgf("Failed to get cart by customer id %s", customerID)
		return
	}

	// Customer already have a shopping cart so we just need to create or update cart item based on condition
	if cart.ID > 0 {
		cartItem, errCartItem := s.CartItemRepo.GetByCartIDAndProductID(ctx.Request().Context(), cart.ID, product.ID)
		if errCartItem != nil && !errors.Is(errCartItem, gorm.ErrRecordNotFound) { // if cart item not found, then create a new one
			log.Err(err).Msgf("Failed to get cart item by cart id %s", cart.ID)
			return
		}

		// There is already the same product in the cart
		if cartItem.ID > 0 {
			cartItem.Quantity += req.Quantity
			cartItem.UpdatedAt = time.Now()
			err = s.CartItemRepo.Update(ctx.Request().Context(), cartItem)
			if err != nil {
				log.Err(err).Msgf("Failed to update cart item for id %s", cartItem.ID)
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
				log.Err(err).Msgf("Failed to create a new cart item for cart id %s", cart.ID)
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
		log.Err(err).Msgf("Failed to create a new cart item for cart id %s", newCart.ID)
		return
	}

	commitRes := trx.Commit()
	if commitRes.Error != nil {
		log.Err(err).Msg("Failed to commit transaction for add product flow")
		return
	}

	return
}
