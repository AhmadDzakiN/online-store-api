package model

import "time"

type OrderItem struct {
	ID        uint64    `gorm:"column:id;primary_key;auto_increment"`
	OrderID   string    `gorm:"column:order_id;not null"`
	ProductID string    `gorm:"column:product_id;not null"`
	Quantity  uint      `gorm:"column:quantity;not null"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}

func (oi *OrderItem) TableName() string {
	return "order_items"
}
