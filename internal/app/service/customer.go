package service

type CustomerService struct {
}

type ICustomerService interface {
}

func NewCustomerService() ICustomerService {
	return &CustomerService{}
}
