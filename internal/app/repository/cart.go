package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"online-store-api/internal/app/model"
)

type CartRepository struct {
	db *gorm.DB
}

type ICartRepository interface {
	GetActiveByCustomerID(ctx context.Context, customerID string) (cart model.Cart, err error)
	Create(ctx context.Context, newCart *model.Cart, trx *gorm.DB) (err error)
}

func NewCartRepository(db *gorm.DB) ICartRepository {
	return &CartRepository{
		db: db,
	}
}

func (r *CartRepository) GetActiveByCustomerID(ctx context.Context, customerID string) (cart model.Cart, err error) {
	result := r.db.WithContext(ctx).First(&cart, "customer_id = ? AND deleted_at IS NULL", customerID)
	if result.Error != nil {
		err = result.Error
		return
	}

	return
}

func (r *CartRepository) Create(ctx context.Context, newCart *model.Cart, trx *gorm.DB) (err error) {
	var db *gorm.DB
	if trx == nil {
		db = r.db.WithContext(ctx)
	} else {
		db = trx
	}

	result := db.Create(newCart)
	if result.Error != nil {
		err = result.Error
		return
	}

	if result.RowsAffected < 1 {
		err = errors.New("no new cart data is created")
		return
	}

	return
}

func (r *CartRepository) Update(ctx context.Context, updatedCart model.Cart) (err error) {
	result := r.db.WithContext(ctx).Save(updatedCart)
	if result.Error != nil {
		err = result.Error
		return
	}

	if result.RowsAffected < 1 {
		err = errors.New("no new cart data is updated")
		return
	}

	return
}
