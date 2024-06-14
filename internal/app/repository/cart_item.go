package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"online-store-api/internal/app/model"
)

type CartItemRepository struct {
	db *gorm.DB
}

type ICartItemRepository interface {
	GetByCartIDAndProductID(ctx context.Context, cartID uint64, productID string) (cartItem model.CartItem, err error)
	Create(ctx context.Context, newCartItem *model.CartItem, trx *gorm.DB) (err error)
	Update(ctx context.Context, updatedCartItem model.CartItem) (err error)
}

func NewCartItemRepository(db *gorm.DB) ICartItemRepository {
	return &CartItemRepository{
		db: db,
	}
}

func (r *CartItemRepository) GetByCartIDAndProductID(ctx context.Context, cartID uint64, productID string) (cartItem model.CartItem, err error) {
	result := r.db.WithContext(ctx).First(&cartItem, "cart_id = ? AND product_id = ?", cartID, productID)
	if result.Error != nil {
		err = result.Error
		return
	}

	return
}

func (r *CartItemRepository) Create(ctx context.Context, newCartItem *model.CartItem, trx *gorm.DB) (err error) {
	var db *gorm.DB
	if trx == nil {
		db = r.db.WithContext(ctx)
	} else {
		db = trx
	}

	result := db.Create(newCartItem)
	if result.Error != nil {
		err = result.Error
		return
	}

	if result.RowsAffected < 1 {
		err = errors.New("no new cart item data is created")
		return
	}

	return
}

func (r *CartItemRepository) Update(ctx context.Context, updatedCartItem model.CartItem) (err error) {
	result := r.db.WithContext(ctx).Updates(updatedCartItem)
	if result.Error != nil {
		err = result.Error
		return
	}

	if result.RowsAffected < 1 {
		err = errors.New("no new cart item data is updated")
		return
	}

	return
}