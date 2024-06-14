package service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"online-store-api/internal/app/model"
	"online-store-api/internal/app/payloads"
	"online-store-api/internal/app/repository"
	"online-store-api/internal/pkg/hash"
)

type CustomerService struct {
	Validate     *validator.Validate
	CustomerRepo repository.ICustomerRepository
}

type ICustomerService interface {
	Register(ctx echo.Context, req payloads.RegisterRequest) (err error)
}

func NewCustomerService(validate *validator.Validate, customerRepo repository.ICustomerRepository) ICustomerService {
	return &CustomerService{
		Validate:     validate,
		CustomerRepo: customerRepo,
	}
}

func (s *CustomerService) Register(ctx echo.Context, req payloads.RegisterRequest) (err error) {
	err = s.Validate.Struct(req)
	if err != nil {
		log.Err(err).Msg("Invalid register request body")
		err = echo.NewHTTPError(http.StatusBadRequest, "Invalid or empty register body request")
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
		log.Err(err).Msgf("Failed to create a new Customer for email %s", req.Email)
		return
	}

	return
}
