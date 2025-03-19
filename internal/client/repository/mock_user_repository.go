package repository

import (
	"1mao/internal/client/domain"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
  mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.Client, error){
	args := m.Called(email)
	if args.Get(0) == nil{
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Client), args.Error(1)
}

func (m *MockUserRepository) Create(user *domain.Client) error{
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(userID uint) (*domain.Client, error){
	args := m.Called(userID)
	if args.Get(0) == nil{
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Client), args.Error(1)
}

func (m *MockUserRepository) GetAllUsers()([]domain.Client, error){
	args := m.Called()
	return args.Get(0).([]domain.Client), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *domain.Client) error{
	args := m.Called(user)
	return args.Error(0)
}
