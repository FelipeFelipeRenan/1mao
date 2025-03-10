package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint
	Email    string
	Password string
	Role     string
}

// UserRepository define os metodos necessarios para autenticação
type UserRepository interface {
	FindByEmail(email string) (*User, error)
}

// AuthService define o serviço de autenticação
type AuthService interface {
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo UserRepository
}

func NewAuthService(userRepo UserRepository) AuthService{
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(email, password string)(string, error){
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("usuario ou senha invalidos")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil{
		return "", errors.New("usuario ou senha invalidos")
	}

	// criar token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role": user.Role,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	secretKey := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}