package service

import (
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
	"time"
)

type OrderService struct {
	Validate      *validator.Validate
	DB            *gorm.DB
	OrderRepo     repository.IOrderRepository
	OrderItemRepo repository.IOrderItemRepository
	ProductRepo   repository.IProductRepository
	CartRepo      repository.ICartRepository
	CartItemRepo  repository.ICartItemRepository
}

type IOrderService interface {
	Checkout(ctx echo.Context, req payloads.CheckoutRequest) (resp payloads.CheckoutResponse, err error)
}

func NewOrderService(validate *validator.Validate, db *gorm.DB, orderRepo repository.IOrderRepository, orderItemRepo repository.IOrderItemRepository,
	productRepo repository.IProductRepository, cartRepo repository.ICartRepository, cartItemRepo repository.ICartItemRepository) IOrderService {
	return &OrderService{
		Validate:      validate,
		DB:            db,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		ProductRepo:   productRepo,
		CartRepo:      cartRepo,
		CartItemRepo:  cartItemRepo,
	}
}

func (s *OrderService) Checkout(ctx echo.Context, req payloads.CheckoutRequest) (resp payloads.CheckoutResponse, err error) {
	err = s.Validate.Struct(req)
	if err != nil {
		log.Err(err).Msg("Invalid checkout request body")
		err = echo.NewHTTPError(http.StatusBadRequest, "invalid checkout request body")
		return
	}

	jwtClaims, err := jwt.GetTokenClaims(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to get jwt token claims")
		err = echo.NewHTTPError(http.StatusUnauthorized, "Failed to get jwt token claims")
		return
	}

	customerID := jwtClaims["customer_id"].(string)
	cart, err := s.CartRepo.GetActiveByIDAndCustomerID(ctx.Request().Context(), req.CartID, customerID)
	if err != nil {
		log.Err(err).Msgf("Failed to get cart by custoemr id %s and cart id %d", customerID, req.CartID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = echo.NewHTTPError(http.StatusNotFound, "customer cart is not found")
			return
		}
		return
	}

	var (
		productIDs  []string
		cartItemIDs []uint64
		totalAmount uint64
		orderItems  []model.OrderItem
	)

	for _, checkOutItem := range req.Items {
		productIDs = append(productIDs, checkOutItem.ProductID)
		cartItemIDs = append(cartItemIDs, checkOutItem.CartItemID)
		totalAmount += checkOutItem.Price * uint64(checkOutItem.Quantity)

		orderItems = append(orderItems, model.OrderItem{
			ProductID: checkOutItem.ProductID,
			Quantity:  checkOutItem.Quantity,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		})
	}

	// Verify all product ids input are exists in database
	err = s.ProductRepo.CheckExistingIDs(ctx.Request().Context(), productIDs)
	if err != nil {
		log.Err(err).Msgf("There are products that do not exist for customer id %s", customerID)
		err = echo.NewHTTPError(http.StatusNotFound, "some products do not exists")
		return
	}

	// Verify all cart item ids input are exists in database
	err = s.CartItemRepo.CheckExistingIDs(ctx.Request().Context(), cartItemIDs, cart.ID)
	if err != nil {
		log.Err(err).Msgf("There are chart items that do not exist for customer id %s", customerID)
		err = echo.NewHTTPError(http.StatusNotFound, "some cart items do not exists")
		return
	}

	// Start db trx
	trx := s.DB.Begin().WithContext(ctx.Request().Context())
	defer trx.Rollback()

	newOrder := model.Order{
		CustomerID:  customerID,
		TotalAmount: totalAmount,
		StatusID:    constants.OrderStatusPending.EnumIndex(),
	}

	// Create order
	err = s.OrderRepo.Create(ctx.Request().Context(), &newOrder, trx)
	if err != nil {
		log.Err(err).Msgf("Failed to create a new order for customer id %s", customerID)
		return
	}

	for i := range orderItems {
		orderItems[i].OrderID = newOrder.ID
	}

	// Create order items with batching feature
	err = s.OrderItemRepo.BatchCreate(ctx.Request().Context(), orderItems, trx)
	if err != nil {
		log.Err(err).Msgf("Failed to batch create order items for customer id %s", customerID)
		return
	}

	// Update the cart of items that have been successfully ordered to fill the deleted_at column
	// If all cart items are ordered, We do not need to update cart data too. Just let it empty so it can be use again
	err = s.CartItemRepo.UpdateToDeleted(ctx.Request().Context(), cartItemIDs, trx)
	if err != nil {
		log.Err(err).Msgf("Failed to update cart items to be deleted for customer id %s", customerID)
		return
	}

	commitRes := trx.Commit()
	if commitRes.Error != nil {
		log.Err(err).Msg("Failed to commit transaction for checkout flow")
		return
	}

	resp = payloads.CheckoutResponse{
		OrderID:     newOrder.ID,
		Status:      constants.OrderStatusPending.String(),
		TotalAmount: totalAmount,
	}

	return
}
