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
	//ProductHandler  handler.ProductHandler
	Config *viper.Viper
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

	customerGroup := e.Group("/customers")
	{
		customerGroup.POST("/register", routeCfg.CustomerHandler.Register)
		customerGroup.POST("/login", routeCfg.CustomerHandler.Login)
	}

	cartGroup := e.Group("/carts")
	{
		cartGroup.Use(echojwt.JWT([]byte(routeCfg.Config.GetString("JWT_SECRET_KEY"))))
		cartGroup.POST("", routeCfg.CartHandler.AddProduct)
	}
}
