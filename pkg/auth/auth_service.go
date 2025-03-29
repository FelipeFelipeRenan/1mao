package auth

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// User representa um usuário normal
type User struct {
	ID       uint
	Email    string
	Password string
}

// Professional representa um profissional
type Professional struct {
	ID       uint
	Email    string
	Password string
}

type Claims struct {
    UserID uint   `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}
// Repositórios de usuários e profissionais
type UserRepository interface {
	FindByEmail(email string) (*User, error)
}

type ProfessionalRepository interface {
	FindByEmail(email string) (*Professional, error)
}

// AuthService define o serviço de autenticação para ambos
type AuthService interface {
	Login(email, password string) (string, error)
	FindByEmail(email string) (*User, error)
}

type authService struct {
	userRepo         UserRepository
	professionalRepo ProfessionalRepository
}

// 🔹 Construtor do AuthService
func NewAuthService(userRepo UserRepository, professionalRepo ProfessionalRepository) AuthService {
	return &authService{userRepo: userRepo, professionalRepo: professionalRepo}
}

// 🔹 Método de login corrigido
func (s *authService) Login(email, password string) (string, error) {
	log.Println("🔵 Tentando login para:", email)

	var userID uint
	var hashedPassword string
	var role string

	// 🔍 Primeiro, tentamos buscar o usuário
	if s.userRepo != nil {
		user, err := s.userRepo.FindByEmail(email)
		if err == nil {
			log.Println("🟢 Usuário encontrado:", user.Email)
			userID = user.ID
			hashedPassword = user.Password
			role = "user"
		} else {
			log.Println("⚠️ Usuário não encontrado no repositório de usuários.")
		}
	}

	// 🔍 Se não encontramos um usuário, tentamos buscar um profissional
	if role == "" && s.professionalRepo != nil {
		professional, err := s.professionalRepo.FindByEmail(email)
		if err == nil {
			log.Println("🟢 Profissional encontrado:", professional.Email)
			userID = professional.ID
			hashedPassword = professional.Password
			role = "professional"
		} else {
			log.Println("⚠️ Profissional não encontrado no repositório de profissionais.")
		}
	}

	// Se nenhum usuário ou profissional foi encontrado, retorna erro
	if role == "" {
		log.Println("❌ Nenhum usuário ou profissional encontrado para o e-mail:", email)
		return "", errors.New("usuário ou senha inválidos")
	}

	// 🔑 Verificar a senha
	log.Println("🟡 Comparando senha fornecida com hash armazenado.")
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		log.Println("❌ Senha incorreta para:", email)
		return "", errors.New("usuário ou senha inválidos")
	}

	log.Println("✅ Senha correta! Gerando token JWT...")

	// 🔐 Criar o token JWT
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Println("⚠️ Chave secreta JWT não está configurada!")
		return "", errors.New("erro interno na autenticação")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("❌ Erro ao gerar token JWT:", err)
		return "", errors.New("erro ao gerar token de autenticação")
	}

	log.Println("✅ Token JWT gerado com sucesso para:", email)
	return tokenString, nil
}

func (s *authService) FindByEmail(email string) (*User, error) {
	// Primeiro, tenta encontrar um usuário normal
	if s.userRepo != nil {
		user, err := s.userRepo.FindByEmail(email)
		if err == nil {
			return user, nil
		}
	}

	// Se não encontrou no userRepo, tenta no professionalRepo
	if s.professionalRepo != nil {
		professional, err := s.professionalRepo.FindByEmail(email)
		if err == nil {
			return &User{
				ID:       professional.ID,
				Email:    professional.Email,
				Password: professional.Password,
			}, nil
		}
	}

	return nil, errors.New("usuário não encontrado")
}
