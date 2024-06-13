package repository

import "gorm.io/gorm"

type ProductCategoryRepository struct {
	db *gorm.DB
}

type IProductCategoryRepository interface {
}

func NewProductCategoryRepository(db *gorm.DB) IProductCategoryRepository {
	return &ProductCategoryRepository{
		db: db,
	}
}
