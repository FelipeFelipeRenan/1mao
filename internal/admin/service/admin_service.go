package service

import (
	"1mao/internal/admin/repository"
	"1mao/pkg/auth"
)

// 🔹 Definição correta do AdminService
type AdminService struct {
	repo    repository.AdminRepository // ❗ Sem o ponteiro, assumindo que é uma interface
	authSvc auth.AuthService
}

// 🔹 Adapter para o AuthService, convertendo Admin para User
type adminServiceAdapter struct {
	repo repository.AdminRepository
}

func (a *adminServiceAdapter) FindByEmail(email string) (*auth.User, error) {
	admin, err := a.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return &auth.User{
		ID:       admin.ID,
		Email:    admin.Email,
		Password: admin.Password,
	}, nil
}

// 🔹 Função para criar o AdminService corretamente
func NewAdminService(repo repository.AdminRepository) *AdminService {
	authRepo := &adminServiceAdapter{repo: repo}
	authSvc := auth.NewAuthService(authRepo, nil) // 🔹 Passando 'nil' para o ProfessionalRepository

	return &AdminService{repo: repo, authSvc: authSvc}
}
