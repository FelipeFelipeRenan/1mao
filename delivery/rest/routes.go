package routes

import (
	"1mao/delivery/rest/routes"
	"1mao/internal/middleware"
	"1mao/internal/client/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// SetupRoutes configura todas as rotas do sistema
func SetupRoutes(db *gorm.DB, clientService *service.ClientService) *mux.Router {
	router := mux.NewRouter()

	// Middlewares globais
	router.Use(middleware.LoggerMiddleware)
	router.Use(middleware.RateLimitMiddleware)
	router.Use(middleware.CircuitBreakerMiddleware)

	// Rota de health check
	routes.HealthRoutes(router, db)	
	// Rota de notificação
	routes.RegisterNotificationRoutes(router)
	// Rota de profissionais
	routes.ProfessionalRoutes(router, db)
	// Rotas de usuário (autenticação e CRUD)
	routes.UserRoutes(router, clientService)

	return router
}

