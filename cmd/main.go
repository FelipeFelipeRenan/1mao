package main

import (
	"1mao/internal/user/delivery/httpa"
	"1mao/internal/user/domain"
	"1mao/internal/user/repository"
	"1mao/internal/user/service"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// Configurações de conexão de banco de dados
	dsn := "user=postgres password=123456 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Taipei"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar no banco de dados: ", err)
	}

	// Migrar tabelas
	db.AutoMigrate(&domain.User{})

	// Instanciar os serviçoes
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	userHandler := httpa.NewUserHandler(authService)

	// Configuração de Router
	router := mux.NewRouter()
	router.HandleFunc("/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")

	os.Setenv("JWT_SECRET", "teste")

	fmt.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}