package model

import (
	"database/sql"
	"time"
)

type CartItem struct {
	ID        uint64       `gorm:"column:id;primary_key;autoIncrement"`
	CartID    uint64       `gorm:"column:cart_id;not null"`
	ProductID string       `gorm:"column:product_id;not null"`
	Quantity  uint         `gorm:"column:quantity;not null"`
	DeletedAt sql.NullTime `gorm:"column:deleted_at"`
	CreatedAt time.Time    `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time    `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}

type CartItemJoinProduct struct {
	ProductID       string    `gorm:"column:product_id"`
	ProductName     string    `gorm:"column:name"`
	ProductPrice    uint64    `gorm:"column:price"`
	ProductQuantity uint      `gorm:"column:quantity"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (ci *CartItem) TableName() string {
	return "cart_items"
}
