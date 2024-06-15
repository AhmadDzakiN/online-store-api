package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"online-store-api/internal/app/constants"
	"online-store-api/internal/app/model"
)

type OrderItemRepository struct {
	db *gorm.DB
}

type IOrderItemRepository interface {
	BatchCreate(ctx context.Context, newOrderItems []model.OrderItem, trx *gorm.DB) (err error)
}

func NewOrderItemRepository(db *gorm.DB) IOrderItemRepository {
	return &OrderItemRepository{
		db: db,
	}
}

func (r *OrderItemRepository) BatchCreate(ctx context.Context, newOrderItems []model.OrderItem, trx *gorm.DB) (err error) {
	var db *gorm.DB
	if trx == nil {
		db = r.db.WithContext(ctx)
	} else {
		db = trx
	}

	result := db.CreateInBatches(newOrderItems, constants.LimitInsertBatch)
	if result.Error != nil {
		err = result.Error
		return
	}

	if result.RowsAffected != int64(len(newOrderItems)) {
		err = fmt.Errorf("expected %d order items created. but got %d", len(newOrderItems), result.RowsAffected)
		return
	}

	return
}
