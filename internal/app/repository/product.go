package repository

import (
	"context"
	"gorm.io/gorm"
	"online-store-api/internal/app/constants"
	"online-store-api/internal/app/model"
	"time"
)

type ProductRepository struct {
	db *gorm.DB
}

type IProductRepository interface {
	GetByID(ctx context.Context, productID string) (product model.Product, err error)
	GetListByCategoryID(ctx context.Context, categoryID string, lastCreated int64) (products []model.Product, err error)
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

func (r *ProductRepository) GetListByCategoryID(ctx context.Context, categoryID string, lastCreated int64) (products []model.Product, err error) {
	query := r.db.WithContext(ctx).Select("id, name, price, updated_at").Where("category_id = ?", categoryID).Order("updated_at DESC")

	if lastCreated > 0 {
		query = query.Where("updated_at < ?", time.Unix(lastCreated, 0))
	}

	query.Limit(constants.LimitDataPerPage)
	result := query.Find(&products)
	if result.Error != nil {
		err = result.Error
		return
	}

	return
}
