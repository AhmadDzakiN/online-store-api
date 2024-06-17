package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"online-store-api/internal/app/cache"
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
	Cache     *redis.Client
}

func BootstrapApp(config *BootstrapAppConfig) {
	cache := cache.NewCache(config.Cache)

	cartRepository := repository.NewCartRepository(config.DB)
	cartItemRepository := repository.NewCartItemRepository(config.DB)
	customerRepository := repository.NewCustomerRepository(config.DB)
	orderRepository := repository.NewOrderRepository(config.DB)
	orderItemRepository := repository.NewOrderItemRepository(config.DB)
	productRepository := repository.NewProductRepository(config.DB)

	cartService := service.NewCartService(config.Validator, config.DB, cache, config.Config, cartRepository, cartItemRepository, productRepository)
	customerService := service.NewCustomerService(config.Validator, config.Config, customerRepository, cartRepository)
	orderService := service.NewOrderService(config.Validator, config.DB, cache, config.Config, orderRepository, orderItemRepository, productRepository, cartRepository, cartItemRepository)
	productService := service.NewProductService(config.Validator, config.Config, productRepository, cache)

	cartHandler := handler.NewCartHandler(cartService)
	customerHandler := handler.NewCustomerHandler(customerService)
	orderHandler := handler.NewOrderHandler(orderService)
	productHandler := handler.NewProductHandler(productService)

	routeCfg := router.RouteConfig{
		CartHandler:     cartHandler,
		CustomerHandler: customerHandler,
		OrderHandler:    orderHandler,
		ProductHandler:  productHandler,
		Config:          config.Config,
	}

	router.NewRouter(routeCfg, config.Echo)

	return
}
