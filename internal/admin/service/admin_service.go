package service

import (
	"1mao/internal/admin/repository"
	"1mao/pkg/auth"
)

type AdminService struct {
	repo *repository.AdminRepository
	authSvc auth.AuthService
}

type adminServiceAdapter struct {
	repo *repository.AdminRepository
}

func (a *adminServiceAdapter) FindByEmail(email string) (*auth.User, error){
	admin, err := a.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return &auth.User{
		ID: admin.ID,
		Email: admin.Email,
		Password: admin.Password,
		Role: "admin",
	}, nil
}

func NewAdminRepository(repo *repository.AdminRepository) *AdminService{
	authRepo := &adminServiceAdapter{repo: repo}
	authSvc := auth.NewAuthService(authRepo)
	
	return &AdminService{repo: repo, authSvc: authSvc}
	
}
