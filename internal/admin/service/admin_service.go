package service

import (
	"1mao/internal/admin/domain"
	"1mao/internal/admin/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	repo *repository.AdminRepository
}

func NewAdminRepository(repo *repository.AdminRepository) *AdminService{
	return &AdminService{repo: repo}
}

func (s *AdminService) RegisterAdmin(name, email, password string) error{
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &domain.AdminUser{
		Name: name,
		Email: email,
		Password: string(hashedPassword),
	}
	return s.repo.Create(admin)
}

func (s *AdminService) Login(email, password string) (*domain.AdminUser, error){
	admin, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("admin nao encotrado")
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return nil, errors.New("senha invalida")
	}
	return admin, nil
}