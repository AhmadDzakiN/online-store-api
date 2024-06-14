package config

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"online-store-api/internal/pkg/builder"
)

func SetCustomErrorHandler(e *echo.Echo) {
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		var response builder.ErrorResponse
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			response = builder.BuildErrorResponse(he.Message.(string))
		} else {
			code = code
			response = builder.BuildErrorResponse("Oops, something wrong with the server. Please try again later")
		}
		// Send JSON response
		if !c.Response().Committed {
			c.JSON(code, response)
		}
	}
}
