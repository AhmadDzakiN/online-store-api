package handler

import "online-store-api/internal/app/service"

type CartHandler struct {
	CartService service.ICartService
}

func NewCartHandler(cartService service.ICartService) *CartHandler {
	return &CartHandler{CartService: cartService}
}
