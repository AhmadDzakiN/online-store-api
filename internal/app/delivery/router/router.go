package router

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"net/http"
	"online-store-api/internal/app/delivery/handler"
)

type RouteConfig struct {
	CustomerHandler *handler.CustomerHandler
	CartHandler     *handler.CartHandler
	OrderHandler    *handler.OrderHandler
	ProductHandler  *handler.ProductHandler
	Config          *viper.Viper
}

func NewRouter(routeCfg RouteConfig, e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			return true, nil
		},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// health check
	e.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "OK!")
	})

	jwtSecretKey := routeCfg.Config.GetString("JWT_SECRET_KEY")

	customerGroup := e.Group("/customers")
	{
		customerGroup.POST("/register", routeCfg.CustomerHandler.Register)
		customerGroup.POST("/login", routeCfg.CustomerHandler.Login)
	}

	cartGroup := e.Group("/carts")
	{
		cartGroup.Use(echojwt.JWT([]byte(jwtSecretKey)))
		cartGroup.POST("", routeCfg.CartHandler.AddProduct)
		cartGroup.DELETE("/:product_id", routeCfg.CartHandler.DeleteProduct)
		cartGroup.GET("", routeCfg.CartHandler.View)
	}

	productGroup := e.Group("/products")
	{
		productGroup.GET("/categories/:category_id", routeCfg.ProductHandler.ViewByCategoryID)
	}

	orderGroup := e.Group("/orders")
	{
		orderGroup.Use(echojwt.JWT([]byte(jwtSecretKey)))
		orderGroup.POST("", routeCfg.OrderHandler.Checkout)
	}
}
