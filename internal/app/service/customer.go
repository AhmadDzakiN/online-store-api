package service

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
	"online-store-api/internal/app/model"
	"online-store-api/internal/app/payloads"
	"online-store-api/internal/app/repository"
	"online-store-api/internal/pkg/hash"
	"online-store-api/internal/pkg/jwt"
)

type CustomerService struct {
	Validate     *validator.Validate
	Config       *viper.Viper
	CustomerRepo repository.ICustomerRepository
	CartRepo     repository.ICartRepository
}

type ICustomerService interface {
	Register(ctx echo.Context, req payloads.RegisterRequest) (err error)
	Login(ctx echo.Context, req payloads.LoginRequest) (resp payloads.LoginResponse, err error)
}

func NewCustomerService(validate *validator.Validate, config *viper.Viper, customerRepo repository.ICustomerRepository,
	cartRepo repository.ICartRepository) ICustomerService {
	return &CustomerService{
		Validate:     validate,
		Config:       config,
		CustomerRepo: customerRepo,
		CartRepo:     cartRepo,
	}
}

func (s *CustomerService) Register(ctx echo.Context, req payloads.RegisterRequest) (err error) {
	err = s.Validate.Struct(req)
	if err != nil {
		log.Err(err).Msg("Invalid register request body")
		err = echo.NewHTTPError(http.StatusBadRequest, "invalid or empty register body request")
		return
	}

	customer, err := s.CustomerRepo.GetByEmail(ctx.Request().Context(), req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Err(err).Msgf("Failed to get customer for email %s", req.Email)
		return
	}

	// There is a customer exist with the email provided
	if customer.ID != "" {
		log.Warn().Msgf("Customer with email %s already exists", req.Email)
		err = echo.NewHTTPError(http.StatusConflict, "customer with email given is already exists")
		return
	}

	// Hash the input password
	hashedPass, err := hash.HashPassword(req.Password)
	if err != nil {
		log.Err(err).Msg("Failed to hash password")
		return
	}

	newCustomer := model.Customer{
		Email:    req.Email,
		Password: hashedPass,
		Name:     req.Name,
		Address:  req.Address,
	}

	err = s.CustomerRepo.Create(ctx.Request().Context(), &newCustomer)
	if err != nil {
		log.Err(err).Msgf("Failed to create a new customer for email %s", req.Email)
		return
	}

	// Create a new cart asynchronously, so it can speed up add product to cart for the first time for user
	go func(ctx context.Context, customer model.Customer) {
		err = s.CartRepo.Create(context.Background(), &model.Cart{CustomerID: customer.ID}, nil)
		if err != nil {
			log.Err(err).Msgf("Failed to create a cart for customer %s in register flow", customer.ID)
			return
		}
	}(ctx.Request().Context(), newCustomer)

	return
}

func (s *CustomerService) Login(ctx echo.Context, req payloads.LoginRequest) (resp payloads.LoginResponse, err error) {
	err = s.Validate.Struct(req)
	if err != nil {
		log.Err(err).Msg("Invalid login request body")
		err = echo.NewHTTPError(http.StatusBadRequest, "invalid login request body")
		return
	}

	customer, err := s.CustomerRepo.GetByEmail(ctx.Request().Context(), req.Email)
	if err != nil {
		log.Err(err).Msgf("Failed to get customer for email %s", req.Email)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = echo.NewHTTPError(http.StatusNotFound, "customer is not found")
		}
		return
	}

	isValid := hash.CheckPasswordHash(req.Password, customer.Password)
	if !isValid {
		log.Error().Msgf("Password does not match for customer %s", req.Email)
		err = echo.NewHTTPError(http.StatusUnauthorized, "Password does not match")
		return
	}

	jwtToken, err := jwt.CreateToken(customer.ID, customer.Name, s.Config)
	if err != nil {
		log.Err(err).Msgf("Failed to create token for customer %s", req.Email)
		return
	}

	resp = payloads.LoginResponse{
		Token: jwtToken,
	}

	return
}
