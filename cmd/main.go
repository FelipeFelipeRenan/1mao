package main

import (
	"1mao/config/database"
	routes "1mao/delivery/rest"
	booking "1mao/internal/booking/domain"
	client "1mao/internal/client/domain"
	"1mao/internal/client/repository"
	"1mao/internal/client/service"
	chat "1mao/internal/notification/domain"
	payment "1mao/internal/payment/domain"
	professional "1mao/internal/professional/domain"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// @title		1Mao API
// @version	1.0
func main() {

	// Carregar variáveis de ambiente
	godotenv.Load(".env")
	dbConfig := database.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("falha na inicialização do banco de dados: %v", err)
	}

	// Conectar ao banco de dados
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	models := []interface{}{
		&client.Client{},
		&professional.Professional{},
		&chat.Message{},
		&booking.Booking{},
		&booking.Availability{},
		&payment.Transaction{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil{
			log.Fatalf("erro ao migrar modelo: %v", err)
		}
		log.Printf("tabela para %T criada com sucesso", model)
	}

	// Instanciar serviços
	userRepo := repository.NewUserRepository(db)
	clientService := service.NewClientService(userRepo)

	// Configuração de rotas
	router := routes.SetupRoutes(db, &clientService)

	// Definir JWT_SECRET na variável de ambiente
	token := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", token)

	// Obter porta da aplicação
	server_port := os.Getenv("APP_PORT")

	fmt.Printf("---- Servidor rodando na porta %s\n ----", server_port)
	log.Fatal(http.ListenAndServe(":"+server_port, router))
}
