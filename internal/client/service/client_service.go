package service

import (
	"1mao/internal/client/domain"
	"1mao/internal/client/repository"
	"1mao/pkg/auth"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ClientService interface {
	Register(user *domain.Client) error
	FindByEmail(email string)(*auth.User, error)
	GetUserByID(userID uint) (*domain.Client, error)
	Login(email, password string) (string, error)
	GetAllUsers() ([]domain.Client, error)
	ForgotPassword(email string) (string, error)
}

type clientService struct {
	userRepo repository.UserRepository
	authSvc auth.AuthService
}

type clientAuthAdapter struct {
	repo repository.UserRepository
}

// Implementação para `FindByEmail`
func (a *clientAuthAdapter) FindByEmail(email string) (*auth.User, error) {
	user, err := a.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return &auth.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (s *clientService) FindByEmail(email string) (*auth.User, error) {
	return s.authSvc.FindByEmail(email) // 🔹 Agora redireciona corretamente
}



func NewClientService(userRepo repository.UserRepository) ClientService {
	authRepo := &clientAuthAdapter{repo: userRepo} // 🔹 Criamos o adapter
	authSvc := auth.NewAuthService(authRepo, nil) // 🔹 Agora passamos o adapter para AuthService

	return &clientService{
		userRepo: userRepo,
		authSvc:  authSvc,
	}
}

// Adapter para conectar UserRepository ao AuthService


func (s *clientService) Register(user *domain.Client) error {
	// Hash da senha do usuario
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	// Salvar no banco
	return s.userRepo.Create(user)
}

func (s *clientService) GetUserByID(userID uint) (*domain.Client, error) {
	return s.userRepo.FindByID(userID)
}

func (s *clientService) GetAllUsers() ([]domain.Client, error) {
	return s.userRepo.GetAllUsers()
}

func (s *clientService) Login(email, password string)(string, error){
	return s.authSvc.Login(email, password)
}
func (s *clientService) ForgotPassword(email string) (string, error){
	
	user, err := s.userRepo.FindByEmail(email)
	if err != nil{
		return "", errors.New("usuario nao encontrado")
	}

	// Gerar um token unico para redefinir a senha
	token := uuid.New().String()
	user.ResetToken = token
	user.ResetTokenExpiry = time.Now().Add(1 * time.Hour) // token expira em 1 hora

	if err := s.userRepo.UpdateUser(user); err != nil{
		return "", errors.New("erro ao salvar token de redefinição")
	}

	if err := sendResetPasswordEmail(user.Email, token); err != nil{
		return "", fmt.Errorf("erro ao enviar email: %w ", err)
	}
	return "Email de recuperação enviado", nil
}