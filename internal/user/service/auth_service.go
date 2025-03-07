package service

import (
	"1mao/internal/user/domain"
	"1mao/internal/user/repository"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(user *domain.User) error
	Login(email, password string) (string, error)
	GetUserByID(userID uint) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
	ForgotPassword(email string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
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

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("usuário ou senha inválidos")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("usuário ou senhas invalidos")
	}

	// criar token jwt

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	secretKey := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func (s *authService) GetUserByID(userID uint) (*domain.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *authService) GetAllUsers() ([]domain.User, error) {
	return s.userRepo.GetAllUsers()
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