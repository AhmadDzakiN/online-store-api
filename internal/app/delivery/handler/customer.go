package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"online-store-api/internal/app/payloads"
	"online-store-api/internal/app/service"
	"online-store-api/internal/pkg/builder"
)

type CustomerHandler struct {
	CustomerService service.ICustomerService
}

func NewCustomerHandler(customerService service.ICustomerService) *CustomerHandler {
	return &CustomerHandler{CustomerService: customerService}
}

func (h *CustomerHandler) Register(ctx echo.Context) (err error) {
	var registerReq payloads.RegisterRequest
	err = ctx.Bind(&registerReq)
	if err != nil {
		log.Err(err).Msg("Invalid register body request")
		err = echo.NewHTTPError(http.StatusBadRequest, "Invalid or empty register body request")
		return
	}

	err = h.CustomerService.Register(ctx, registerReq)
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusCreated, builder.BuildSuccessResponse(nil))
}

func (h *CustomerHandler) Login(ctx echo.Context) (err error) {
	var loginReq payloads.LoginRequest
	err = ctx.Bind(&loginReq)
	if err != nil {
		log.Err(err).Msg("Invalid login body request")
		err = echo.NewHTTPError(http.StatusBadRequest, "Invalid or empty login body request")
		return
	}

	resp, err := h.CustomerService.Login(ctx, loginReq)
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, builder.BuildSuccessResponse(resp))
}
