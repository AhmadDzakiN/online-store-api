package model

import (
	"time"
)

type ProductCategory struct {
	ID          string    `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name        string    `gorm:"column:name;not null"`
	Description string    `gorm:"column:description;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}

func (pc *ProductCategory) TableName() string {
	return "product_categories"
}
