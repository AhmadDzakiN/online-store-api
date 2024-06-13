package repository

import "gorm.io/gorm"

type ProductRepository struct {
	db *gorm.DB
}

type IProductRepository interface {
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{
		db: db,
	}
}
