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
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Carrega variáveis de ambiente do .env
func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️ Aviso: Arquivo .env não encontrado ou não pode ser carregado. Usando variáveis de ambiente padrão.")
	}
}

// Obtém uma variável de ambiente ou retorna um valor padrão
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Conecta ao banco de dados e tenta criá-lo caso não exista
func connectDatabase() *gorm.DB {

	db, err := gorm.Open(postgres.Open("host=db user=postgres password=postgres dbname=1mao port=5432 sslmode=disable"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Printf("❌ Erro ao conectar no banco de dados: %v", err)
	}
	return db
}

func main() {
	// Carregar variáveis de ambiente
	loadEnv()

	// Verificar valores carregados
	log.Println("📌 Configurações carregadas:")
	log.Printf("🔹 DB_HOST: %s", getEnv("DB_HOST", "localhost"))
	log.Printf("🔹 DB_USER: %s", getEnv("DB_USER", "postgres"))
	log.Printf("🔹 DB_NAME: %s", getEnv("DB_NAME", "1mao"))
	log.Printf("🔹 DB_PORT: %s", getEnv("DB_PORT", "5432"))

	// Conectar ao banco de dados
	db := connectDatabase()
	log.Println("✅ Conectado ao banco de dados com sucesso.")


	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Migrar tabelas
	db.Migrator().DropTable(&domain.User{})
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatal("Erro ao migrar modelo", err)
	}
	
	log.Println("Tabela 'user' criada com sucesso")
		
	// Instanciar serviços
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	userHandler := httpa.NewUserHandler(authService)

	// Configuração do Router
	router := mux.NewRouter()
	router.HandleFunc("/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")

	// Definir JWT_SECRET na variável de ambiente
	os.Setenv("JWT_SECRET", getEnv("JWT_SECRET", "defaultt"))

	// Obter porta da aplicação
	port := getEnv("APP_PORT", "8080")

	fmt.Printf("🚀 Servidor rodando na porta %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
