package service

import (
	"1mao/internal/professional/domain"
	"1mao/internal/professional/repository"

	"golang.org/x/crypto/bcrypt"
)

type ProfessionalService interface {
	Register(professional *domain.Professional) error
	GetProfessionalByID(id uint) (*domain.Professional, error)
	GetAllProfessionals() ([]domain.Professional, error)
}

type professionalService struct {
	repo repository.ProfessionalRepository
}

func NewProfessionalService(repo repository.ProfessionalRepository) ProfessionalService {
	return &professionalService{repo: repo}
}

func (s *professionalService) Register(professional *domain.Professional) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(professional.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	professional.Password = string(hashedPassword)
	return s.repo.Create(professional)
}


func (s *professionalService) GetProfessionalByID(id uint)(*domain.Professional, error){
	return s.repo.FindByID(id)
}

func (s *professionalService) GetAllProfessionals()([]domain.Professional, error){
	return s.repo.GetAllProfessionals()
}