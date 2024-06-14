package model

import "time"

type PaymentStatus struct {
	ID          uint      `gorm:"column:id;primary_key;auto_increment"`
	Name        string    `gorm:"column:name;not null"`
	Description string    `gorm:"column:description;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}

func (ps *PaymentStatus) TableName() string {
	return "payment_statuses"
}
