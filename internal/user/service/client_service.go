package service

import (
	"1mao/internal/user/domain"
	"1mao/internal/user/repository"
	"1mao/pkg/auth"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(user *domain.User) error
	FindByEmail(email string)(*auth.User, error)
	GetUserByID(userID uint) (*domain.User, error)
	Login(email, password string) (string, error)
	GetAllUsers() ([]domain.User, error)
	ForgotPassword(email string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	authSvc auth.AuthService
}


func NewAuthService(userRepo repository.UserRepository) AuthService {
	clientRepo := &authService{userRepo: userRepo}
	authSvc := auth.NewAuthService(clientRepo)
	return &authService{userRepo: userRepo, authSvc: authSvc}
}

func (s *authService) FindByEmail(email string)(*auth.User, error){
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return &auth.User{
		ID: user.ID,
		Email: user.Email,
		Password: user.Password,
		Role: string(user.Role),
	}, nil
}

func (s *authService) Register(user *domain.User) error {
	// Hash da senha do usuario
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	// Salvar no banco
	return s.userRepo.Create(user)
}

func (s *authService) GetUserByID(userID uint) (*domain.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *authService) GetAllUsers() ([]domain.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s *authService) Login(email, password string)(string, error){
	return s.authSvc.Login(email, password)
}
func (s *authService) ForgotPassword(email string) (string, error){
	
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