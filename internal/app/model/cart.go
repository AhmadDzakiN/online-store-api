package model

import (
	"database/sql"
	"time"
)

type Cart struct {
	ID         uint64       `gorm:"column:id;primary_key;autoIncrement"`
	CustomerID string       `gorm:"column:customer_id;not null"`
	DeletedAt  sql.NullTime `gorm:"column:deleted_at"`
	CreatedAt  time.Time    `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt  time.Time    `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}

func (c *Cart) TableName() string {
	return "carts"
}
