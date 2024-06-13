package repository

import "gorm.io/gorm"

type CartRepository struct {
	db *gorm.DB
}

type ICartRepository interface {
}

func NewCartRepository(db *gorm.DB) ICartRepository {
	return &CartRepository{
		db: db,
	}
}
