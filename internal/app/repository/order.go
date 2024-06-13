package repository

import "gorm.io/gorm"

type OrderRepository struct {
	db *gorm.DB
}

type IOrderRepository interface {
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{
		db: db,
	}
}
