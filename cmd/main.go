package main

import (
	routes "1mao/delivery/rest"
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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Conecta ao banco de dados e tenta criá-lo caso não exista
func connectDatabase(host string, user string, password string, name string, port string, sslmode string) *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, name, port, sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Printf("❌ Erro ao conectar no banco de dados: %v", err)
	}
	return db
}

// @title		1Mao API
// @version	1.0
func main() {

	// Carregar variáveis de ambiente
	godotenv.Load(".env")

	db_host := os.Getenv("DB_HOST")
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")
	db_sslmode := os.Getenv("DB_SSLMODE")

	// Conectar ao banco de dados
	db := connectDatabase(db_host, db_user, db_password, db_name, db_port, db_sslmode)
	log.Println("✅ Conectado ao banco de dados com sucesso.")

	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	if err := db.AutoMigrate(&client.Client{}); err != nil {
		log.Fatal("Erro ao migrar modelo User:", err)
	}
	log.Println("Tabela 'client' criada com sucesso")

	if err := db.AutoMigrate(&professional.Professional{}); err != nil {
		log.Fatal("Erro ao migrar modelo Professional:", err)
	}
	log.Println("Tabela 'professional' criada com sucesso")

	if err := db.AutoMigrate(&chat.Message{}); err != nil {
		log.Fatal("Erro ao migrar modelo Professional:", err)
	}
	log.Println("Tabela 'message' criada com sucesso")

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
