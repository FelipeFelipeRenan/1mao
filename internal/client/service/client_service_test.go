package service

import (
	"1mao/internal/client/domain"
	"1mao/internal/client/repository"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)


func hashPassword(password string) string{
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

func TestRegister_Sucess(t *testing.T){
	mockRepo := new(repository.MockClientRepository)
	clientService := NewClientService(mockRepo)

	user := &domain.Client{
		Email: "user@email.com",
		Password: "senha123",
		Role: "client",
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.Client")).Return(nil)

	err := clientService.Register(user)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

}

func TestFindByEmail_Success(t *testing.T){
	mockRepo := new(repository.MockClientRepository)
	authService := NewClientService(mockRepo)

	expectedUser := &domain.Client{
		Email: "user@email.com",
		Password: hashPassword("senha123"),
		Role: domain.RoleClient,
	}

	mockRepo.On("FindByEmail", "user@email.com").Return(expectedUser,nil)
	
	user, err := authService.FindByEmail("user@email.com")
	assert.NoError(t, err)

	assert.Equal(t, expectedUser.Email, user.Email)
	mockRepo.AssertExpectations(t)
}

func TestFindbyEmail_NotFound(t *testing.T){
	mockRepo := new(repository.MockClientRepository)
	authService := NewClientService(mockRepo)

	mockRepo.On("FindByEmail", "naoexiste@email.com").Return(nil, errors.New("usuario nao encontrado"))

	user, err := authService.FindByEmail("naoexiste@email.com")
	
	assert.Nil(t, user)
	assert.Equal(t, "usuário não encontrado", err.Error())

	mockRepo.AssertExpectations(t)
}