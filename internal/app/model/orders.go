package model

import "time"

type Order struct {
	ID          string    `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	CustomerID  string    `gorm:"column:customer_id;not null"`
	TotalAmount uint64    `gorm:"column:total_amount;not null"`
	StatusID    uint      `gorm:"column:status_id;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}

func (o *Order) TableName() string {
	return "orders"
}
