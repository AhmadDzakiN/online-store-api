package model

import "time"

type Order struct {
	ID            uint64    `gorm:"primary_key;auto_increment"`
	CustomerID    string    `gorm:"column:customer_id;not null"`
	TotalAmount   uint64    `gorm:"column:total_amount;not null"`
	PaymentStatus uint      `gorm:"column:payment_status;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt     time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}

func (o *Order) TableName() string {
	return "orders"
}
