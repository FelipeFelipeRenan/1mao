package routes

import (
	"1mao/delivery/rest/routes"
	"1mao/internal/client/service"
	"1mao/internal/middleware"
	"1mao/internal/notification/repository"
	"1mao/internal/notification/websocket"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// SetupRoutes configura todas as rotas do sistema
func SetupRoutes(db *gorm.DB, clientService *service.ClientService) *mux.Router {
	router := mux.NewRouter()

	// Criar repositório de mensagens
	messageRepo := repository.NewMessageRepository(db)

	// Criar Hub com repositório de mensagens
	hub := websocket.NewHub(messageRepo)
	go hub.Run()

	// Middlewares globais
	router.Use(middleware.LoggerMiddleware)
	router.Use(middleware.RateLimitMiddleware)
	router.Use(middleware.CircuitBreakerMiddleware)

	// Rota do Swagger
	routes.SwaggerRoutes()
	// Rota de health check
	routes.HealthRoutes(router, db)
	// Rota de notificação
	routes.RegisterNotificationRoutes(router)
	// Rota de chat
	routes.RegisterChatRoutes(router, db, hub)
	// Rota de profissionais
	routes.ProfessionalRoutes(router, db)
	// Rotas de usuário (autenticação e CRUD)
	routes.UserRoutes(router, clientService)

	

	return router
}

