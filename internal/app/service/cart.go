package service

type CartService struct {
}

type ICartService interface {
}

func NewCartService() ICartService {
	return &CartService{}
}
