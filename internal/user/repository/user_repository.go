package repository

import (
	"1mao/internal/user/domain"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(userID uint) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("usuario nao encontrado")
	}
	return &user, nil
}

func (r *userRepository) FindByID(userID uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAllUsers() ([]domain.User, error) {
	var users []domain.User
	result := r.db.Find(&users)
	if result.Error != nil{
		return nil, result.Error
	}
	return users, nil
}
