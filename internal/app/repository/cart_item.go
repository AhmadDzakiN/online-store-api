package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"online-store-api/internal/app/constants"
	"online-store-api/internal/app/model"
	"time"
)

type CartItemRepository struct {
	db *gorm.DB
}

type ICartItemRepository interface {
	GetActiveByCartIDAndProductID(ctx context.Context, cartID uint64, productID string) (cartItem model.CartItem, err error)
	Create(ctx context.Context, newCartItem *model.CartItem, trx *gorm.DB) (err error)
	Update(ctx context.Context, updatedCartItem model.CartItem) (err error)
	GetActiveItemsAndProductsByCartID(ctx context.Context, cartID uint64, lastCreated int64) (data []model.CartItemJoinProduct, err error)
	CheckExistingIDs(ctx context.Context, cartItemIDs []uint64, cartID uint64) (err error)
	UpdateToDeleted(ctx context.Context, cartItemIDs []uint64, trx *gorm.DB) (err error)
}

func NewCartItemRepository(db *gorm.DB) ICartItemRepository {
	return &CartItemRepository{
		db: db,
	}
}

func (r *CartItemRepository) GetActiveByCartIDAndProductID(ctx context.Context, cartID uint64, productID string) (cartItem model.CartItem, err error) {
	result := r.db.WithContext(ctx).First(&cartItem, "cart_id = ? AND product_id = ? AND deleted_at IS NULL AND quantity != 0", cartID, productID)
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
	result := r.db.WithContext(ctx).Save(updatedCartItem)
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

func (r *CartItemRepository) GetActiveByCartID(ctx context.Context, cartID uint64) (cartItem model.CartItem, err error) {
	result := r.db.WithContext(ctx).First(&cartItem, "cart_id = ? AND deleted_at IS NULL AND quantity != 0", cartID)
	if result.Error != nil {
		err = result.Error
		return
	}

	return
}

func (r *CartItemRepository) CountActiveByCartID(ctx context.Context, cartID uint64) (countActiveItems int64, err error) {
	result := r.db.WithContext(ctx).Where("cart_id = ? AND deleted_at IS NULL AND quantity != 0", cartID).Count(&countActiveItems)
	if err != nil {
		err = result.Error
		return
	}

	return
}

func (r *CartItemRepository) GetActiveItemsAndProductsByCartID(ctx context.Context, cartID uint64, lastCreated int64) (data []model.CartItemJoinProduct, err error) {
	query := r.db.WithContext(ctx).Select("c.cart_id, c.id as cart_item_id, c.product_id, p.name, p.price, c.quantity, c.updated_at").Table("cart_items c").
		Joins("INNER JOIN products p ON c.product_id = p.id").Where("c.cart_id = ? AND c.deleted_at IS NULL AND c.quantity != 0", cartID).Order("c.updated_at DESC")

	if lastCreated > 0 {
		query = query.Where("c.updated_at < ?", time.Unix(lastCreated, 0))
	}

	query.Limit(constants.LimitDataPerPage)
	result := query.Find(&data)
	if result.Error != nil {
		err = result.Error
		return
	}

	return
}

func (r *CartItemRepository) CheckExistingIDs(ctx context.Context, cartItemIDs []uint64, cartID uint64) (err error) {
	var validIDs []uint64
	result := r.db.WithContext(ctx).Select("id").Table("cart_items").
		Where("cart_id = ? AND id IN (?) AND deleted_at IS NULL AND quantity != 0", cartID, cartItemIDs).Find(&validIDs)
	if result.Error != nil {
		err = result.Error
		return
	}

	validIDMap := make(map[uint64]bool)
	for _, validID := range validIDs {
		validIDMap[validID] = true
	}

	var invalidIDs []uint64
	for _, id := range cartItemIDs {
		if !validIDMap[id] {
			invalidIDs = append(invalidIDs, id)
		}
	}

	if len(invalidIDs) > 0 {
		err = fmt.Errorf("cart item IDs %v does not exist", invalidIDs)
		return
	}

	return
}

func (r *CartItemRepository) UpdateToDeleted(ctx context.Context, cartItemIDs []uint64, trx *gorm.DB) (err error) {
	var db *gorm.DB
	if trx == nil {
		db = r.db.WithContext(ctx)
	} else {
		db = trx
	}

	now := time.Now()
	result := db.Table("cart_items").Where("id in (?)", cartItemIDs).Updates(map[string]interface{}{
		"deleted_at": sql.NullTime{Time: now, Valid: true},
		"quantity":   0,
		"updated_at": now,
	})
	if result.Error != nil {
		err = result.Error
		return
	}

	if result.RowsAffected != int64(len(cartItemIDs)) {
		err = fmt.Errorf("expected %d cart items updated. but got %d", len(cartItemIDs), result.RowsAffected)
		return
	}

	return
}
