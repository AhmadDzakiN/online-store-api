package service

type ProductService struct {
}

type IProductService interface {
}

func NewProductService() IProductService {
	return &ProductService{}
}
