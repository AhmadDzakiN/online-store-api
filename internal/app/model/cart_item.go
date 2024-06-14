package model

import (
	"time"
)

type CartItem struct {
	ID        uint64    `gorm:"column:id;primary_key;autoIncrement"`
	CartID    uint64    `gorm:"column:cart_id;not null"`
	ProductID string    `gorm:"column:product_id;not null"`
	Quantity  uint      `gorm:"column:quantity;not null"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}

func (ci *CartItem) TableName() string {
	return "cart_items"
}
