package repository

import (
	"1mao/internal/professional/domain"

	"gorm.io/gorm"
)

type ProfessionalRepository interface {
	Create(professional *domain.Professional) error
	FindByID(id uint)(*domain.Professional, error)
	FindByEmail(email string) (*domain.Professional, error)
	GetAllProfessionals()([]domain.Professional, error)

}

type professionalRepository struct {
	db *gorm.DB
}

func NewProfessionalRepository(db *gorm.DB) ProfessionalRepository{
	return &professionalRepository{db: db}
}

func (r *professionalRepository)Create(professional  *domain.Professional)error {
	return r.db.Create(professional).Error
}


func (r *professionalRepository) FindByID(id uint)(*domain.Professional, error){
	var professional domain.Professional
	if err  := r.db.First(&professional, id).Error; err != nil{
		return nil, err
	}
	return &professional, nil
}

func (r *professionalRepository) FindByEmail(email string)(*domain.Professional, error){
	var professional domain.Professional
	if err:= r.db.Where("email = ?", email).First(&professional).Error; err != nil{
		return nil, err
	}
	return &professional, nil
} 


func (r *professionalRepository) GetAllProfessionals()([]domain.Professional, error){
	var professionals []domain.Professional
	if err := r.db.Find(&professionals).Error; err !=nil{
		return nil, err
	}
	return professionals, nil
}
