package repository

import (
	"context"
	"gorm.io/gorm"
	"online-store-api/internal/app/model"
)

type ProductRepository struct {
	db *gorm.DB
}

type IProductRepository interface {
	GetByID(ctx context.Context, productID string) (product model.Product, err error)
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) GetByID(ctx context.Context, productID string) (product model.Product, err error) {
	result := r.db.WithContext(ctx).First(&product, "id = ?", productID)
	if result.Error != nil {
		err = result.Error
		return
	}

	return
}
