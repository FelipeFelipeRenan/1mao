package repository

import "gorm.io/gorm"

type BookingRepository interface {
	
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB)BookingRepository{
	return &bookingRepository{db: db}
}

// TODO: Metodos do repo