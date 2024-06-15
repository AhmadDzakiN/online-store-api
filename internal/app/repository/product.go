package repository

import (
	"context"
	"fmt"
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
	CheckExistingIDs(ctx context.Context, productIDs []string) (err error)
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

func (r *ProductRepository) CheckExistingIDs(ctx context.Context, productIDs []string) (err error) {
	var validIDs []string
	result := r.db.WithContext(ctx).Select("id").Table("products").Where("id IN (?)", productIDs).Find(&validIDs)
	if result.Error != nil {
		err = result.Error
		return
	}

	validIDMap := make(map[string]bool)
	for _, validID := range validIDs {
		validIDMap[validID] = true
	}

	var invalidIDs []string
	for _, id := range productIDs {
		if !validIDMap[id] {
			invalidIDs = append(invalidIDs, id)
		}
	}

	if len(invalidIDs) > 0 {
		err = fmt.Errorf("product IDs %v does not exist", invalidIDs)
		return
	}

	return
}
