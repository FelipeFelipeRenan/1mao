package repository

import (
	"1mao/internal/admin/domain"

	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository{
	return &AdminRepository{db: db}
}

func (r *AdminRepository) Create(admin *domain.AdminUser) error{
	return r.db.Create(admin).Error
}

func (r *AdminRepository) FindByEmail(email string)(*domain.AdminUser, error){
	var admin domain.AdminUser
	err := r.db.Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
