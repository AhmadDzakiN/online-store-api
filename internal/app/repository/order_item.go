package repository

import "gorm.io/gorm"

type OrderItemRepository struct {
	db *gorm.DB
}

type IOrderItemRepository interface {
}

func NewOrderItemRepository(db *gorm.DB) IOrderItemRepository {
	return &OrderItemRepository{
		db: db,
	}
}
