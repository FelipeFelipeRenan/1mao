package repository

import (
	"1mao/internal/client/domain"
	"log"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.Client) error
	FindByEmail(email string) (*domain.Client, error)
	FindByID(userID uint) (*domain.Client, error)
	GetAllUsers() ([]domain.Client, error)
	UpdateUser(user *domain.Client) error 
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.Client) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*domain.	Client, error) {
	var user domain.Client
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		log.Println("⚠️ Erro ao buscar usuário:", result.Error)
		return nil, result.Error
	}
	log.Println("🟢 Usuário encontrado:", user.Email)
	return &user, nil
}




func (r *userRepository) FindByID(userID uint) (*domain.Client, error) {
	var user domain.Client
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAllUsers() ([]domain.Client, error) {
	var users []domain.Client
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *userRepository) UpdateUser(user *domain.Client) error {
	return r.db.Save(user).Error
}
