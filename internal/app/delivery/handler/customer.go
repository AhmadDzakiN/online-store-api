package handler

import (
	"github.com/labstack/echo/v4"
	"online-store-api/internal/app/service"
)

type CustomerHandler struct {
	CustomerService service.ICustomerService
}

func NewCustomerHandler(customerService service.ICustomerService) *CustomerHandler {
	return &CustomerHandler{CustomerService: customerService}
}

func (h *CustomerHandler) Register(ctx echo.Context) (err error) {
	return
}
