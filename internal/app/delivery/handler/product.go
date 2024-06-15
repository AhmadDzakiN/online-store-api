package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"online-store-api/internal/app/service"
	"online-store-api/internal/pkg/builder"
)

type ProductHandler struct {
	ProductService service.IProductService
}

func NewProductHandler(productService service.IProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: productService,
	}
}

func (h *ProductHandler) ViewByCategoryID(ctx echo.Context) (err error) {
	categoryID := ctx.Param("category_id")
	nextToken := ctx.QueryParam("next")
	resp, nextPageToken, err := h.ProductService.ViewByCategoryID(ctx, categoryID, nextToken)
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, builder.BuildSuccessResponse(resp, &nextPageToken))
}
