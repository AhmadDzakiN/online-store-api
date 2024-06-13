package model

import (
	"time"
)

type Cart struct {
	ID         uint64    `gorm:"primary_key;autoIncrement"`
	CustomerID string    `gorm:"column:customer_id;not null"`
	ProductID  string    `gorm:"column:product_id;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
}

func (c *Cart) TableName() string {
	return "carts"
}
