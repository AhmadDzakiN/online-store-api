package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"online-store-api/internal/app/delivery/handler"
	"online-store-api/internal/app/delivery/router"
	"online-store-api/internal/app/repository"
	"online-store-api/internal/app/service"
)

type BootstrapAppConfig struct {
	DB        *gorm.DB
	Validator *validator.Validate
	Config    *viper.Viper
	Echo      *echo.Echo
}

func BootstrapApp(config *BootstrapAppConfig) {
	//cartRepository := repository.NewCartRepository(config.DB)
	customerRepository := repository.NewCustomerRepository(config.DB)
	//orderRepository := repository.NewOrderRepository(config.DB)
	//orderItemRepository := repository.NewOrderItemRepository(config.DB)
	//paymentStatusRepository := repository.NewPaymentStatusRepository(config.DB)
	//productRepository := repository.NewProductRepository(config.DB)
	//productCategoryRepository := repository.NewProductCategoryRepository(config.DB)

	customerService := service.NewCustomerService(config.Validator, customerRepository)

	customerHandler := handler.NewCustomerHandler(customerService)

	routeCfg := router.RouteConfig{
		CustomerHandler: customerHandler,
	}

	router.NewRouter(routeCfg, config.Echo)

	return
}