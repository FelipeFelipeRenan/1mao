package service

import (
	"1mao/internal/professional/domain"
	"1mao/internal/professional/repository"
	"1mao/pkg/auth"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var (
	ctx = context.Background()
)

// 🔹 Interface do serviço de profissionais
type ProfessionalService interface {
	Register(professional *domain.Professional) error
	GetProfessionalByID(id uint) (*domain.Professional, error)
	GetAllProfessionals() ([]domain.Professional, error)
	Login(email, password string) (string, error) // 🔹 Adicionando Login
}

// 🔹 Implementação do serviço de profissionais
type professionalService struct {
	repo     repository.ProfessionalRepository
	authSvc  auth.AuthService
	cache    *redis.Client
	cacheTTL time.Duration
}

// 🔹 Adapter para conectar ProfessionalRepository ao AuthService
type professionalAuthAdapter struct {
	repo repository.ProfessionalRepository
}

// 🔹 Ajustando o retorno para auth.Professional
func (a *professionalAuthAdapter) FindByEmail(email string) (*auth.Professional, error) {
	professional, err := a.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return &auth.Professional{
		ID:       professional.ID,
		Email:    professional.Email,
		Password: professional.Password,
	}, nil
}

// 🔹 Criando o ProfessionalService corretamente
func NewProfessionalService(repo repository.ProfessionalRepository, redisClient *redis.Client) ProfessionalService {
	authRepo := &professionalAuthAdapter{repo: repo}
	authSvc := auth.NewAuthService(nil, authRepo) // 🔹 Passamos nil para UserRepository

	return &professionalService{repo: repo,
		authSvc:  authSvc,
		cache:    redisClient,
		cacheTTL: 30 * time.Minute}
}

// Helper para operações de cache
func (s *professionalService) getFromCache(key string, target interface{}) bool {
	cached, err := s.cache.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	return json.Unmarshal([]byte(cached), target) == nil
}

func (s *professionalService) setCache(key string, value interface{}) error {
	serialized, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.cache.Set(ctx, key, serialized, s.cacheTTL).Err()
}

func (s *professionalService) invalidateCache(pattern string){
	keys, err := s.cache.Keys(ctx, pattern).Result()
	if err == nil {
		s.cache.Del(ctx, keys...)
	}
}

// 🔹 Registro de profissional
func (s *professionalService) Register(professional *domain.Professional) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(professional.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	professional.Password = string(hashedPassword)
	return s.repo.Create(professional)
}

// 🔹 Buscar profissional por ID
func (s *professionalService) GetProfessionalByID(id uint) (*domain.Professional, error) {
	return s.repo.FindByID(id)
}

// 🔹 Buscar todos os profissionais
func (s *professionalService) GetAllProfessionals() ([]domain.Professional, error) {
	return s.repo.GetAllProfessionals()
}

// 🔹 Implementação do Login usando AuthService
func (s *professionalService) Login(email, password string) (string, error) {
	return s.authSvc.Login(email, password)
}
