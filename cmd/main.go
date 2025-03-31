package main

import (
	"1mao/config/database"
	routes "1mao/delivery/rest"
	booking "1mao/internal/booking/domain"
	client "1mao/internal/client/domain"
	"1mao/internal/client/repository"
	"1mao/internal/client/service"
	chat "1mao/internal/notification/domain"
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

	if err := db.AutoMigrate(&client.Client{}); err != nil {
		log.Fatal("Erro ao migrar modelo Client:", err)
	}
	log.Println("Tabela 'client' criada com sucesso")

	if err := db.AutoMigrate(&professional.Professional{}); err != nil {
		log.Fatal("Erro ao migrar modelo Professional:", err)
	}
	log.Println("Tabela 'professional' criada com sucesso")

	if err := db.AutoMigrate(&chat.Message{}); err != nil {
		log.Fatal("Erro ao migrar modelo Message:", err)
	}
	log.Println("Tabela 'message' criada com sucesso")
	if err := db.AutoMigrate(&booking.Booking{}); err != nil {
		log.Fatal("Erro ao migrar modelo Booking:", err)
	}
	log.Println("Tabela 'booking' criada com sucesso")

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

	fmt.Printf("🚀 Servidor rodando na porta %s\n", server_port)
	log.Fatal(http.ListenAndServe(":"+server_port, router))
}
