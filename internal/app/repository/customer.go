package repository

import "gorm.io/gorm"

type CustomerRepository struct {
	db *gorm.DB
}

type ICustomerRepository interface {
}

func NewCustomerRepository(db *gorm.DB) ICustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}
