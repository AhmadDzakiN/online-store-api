package repository

import "gorm.io/gorm"

type PaymentStatusRepository struct {
	db *gorm.DB
}

type IPaymentStatusRepository interface {
}

func NewPaymentStatusRepository(db *gorm.DB) IPaymentStatusRepository {
	return &PaymentStatusRepository{
		db: db,
	}
}
