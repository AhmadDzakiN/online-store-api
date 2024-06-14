package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"online-store-api/internal/app/payloads"
	"online-store-api/internal/app/service"
	"online-store-api/internal/pkg/builder"
)

type CartHandler struct {
	CartService service.ICartService
}

func NewCartHandler(cartService service.ICartService) *CartHandler {
	return &CartHandler{CartService: cartService}
}

func (h *CartHandler) AddProduct(ctx echo.Context) (err error) {
	var addProductReq payloads.AddProductRequest
	err = ctx.Bind(&addProductReq)
	if err != nil {
		log.Err(err).Msg("Invalid add product to cart request body")
		err = echo.NewHTTPError(http.StatusBadRequest, "invalid add product to cart request body")
		return
	}

	err = h.CartService.AddProduct(ctx, addProductReq)
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusCreated, builder.BuildSuccessResponse(nil))
}
