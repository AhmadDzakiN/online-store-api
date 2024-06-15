package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"online-store-api/internal/app/model"
)

type OrderRepository struct {
	db *gorm.DB
}

type IOrderRepository interface {
	Create(ctx context.Context, newOrder *model.Order, trx *gorm.DB) (err error)
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(ctx context.Context, newOrder *model.Order, trx *gorm.DB) (err error) {
	var db *gorm.DB
	if trx == nil {
		db = r.db.WithContext(ctx)
	} else {
		db = trx
	}

	result := db.WithContext(ctx).Create(newOrder)
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
