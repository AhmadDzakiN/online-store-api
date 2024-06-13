package model

import "time"

type Product struct {
	ID          string    `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"column:name;not null"`
	CategoryID  string    `gorm:"column:category_id;not null"`
	Description string    `gorm:"column:description;not null"`
	Price       uint64    `gorm:"column:price;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}

func (p *Product) TableName() string {
	return "products"
}
