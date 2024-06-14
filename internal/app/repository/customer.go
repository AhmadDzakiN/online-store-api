package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"online-store-api/internal/app/model"
)

type CustomerRepository struct {
	db *gorm.DB
}

type ICustomerRepository interface {
	GetByEmail(ctx context.Context, email string) (customer model.Customer, err error)
	Create(ctx context.Context, newGenerator *model.Customer) (err error)
}

func NewCustomerRepository(db *gorm.DB) ICustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (r *CustomerRepository) GetByEmail(ctx context.Context, email string) (customer model.Customer, err error) {
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&customer)
	if result.Error != nil {
		return
	}

	return
}

func (r *CustomerRepository) Create(ctx context.Context, newGenerator *model.Customer) (err error) {
	result := r.db.WithContext(ctx).Create(&newGenerator)
	if result.Error != nil {
		err = result.Error
		return
	}

	if result.RowsAffected < 1 {
		err = errors.New("no new Customer data is created")
		return
	}

	return
}
