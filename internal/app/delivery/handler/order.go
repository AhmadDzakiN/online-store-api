package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"online-store-api/internal/app/payloads"
	"online-store-api/internal/app/service"
	"online-store-api/internal/pkg/builder"
)

type OrderHandler struct {
	OrderService service.IOrderService
}

func NewOrderHandler(orderService service.IOrderService) *OrderHandler {
	return &OrderHandler{
		OrderService: orderService,
	}
}

func (h *OrderHandler) Checkout(ctx echo.Context) (err error) {
	var checkoutListReq payloads.CheckoutRequest
	err = ctx.Bind(&checkoutListReq)
	if err != nil {
		log.Err(err).Msg("Invalid checkout request body")
		err = echo.NewHTTPError(http.StatusBadRequest, "invalid checkout request body")
		return
	}

	resp, err := h.OrderService.Checkout(ctx, checkoutListReq)
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusCreated, builder.BuildSuccessResponse(resp, nil))
}
