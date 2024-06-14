package model

import (
	"time"
)

type Customer struct {
	ID        string    `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Email     string    `gorm:"column:email;unique;not null"`
	Password  string    `gorm:"column:password;not null"`
	Name      string    `gorm:"column:name;not null"`
	Address   string    `gorm:"column:address;not null"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}

func (c *Customer) TableName() string {
	return "customers"
}
