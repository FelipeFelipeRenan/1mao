package routes

import (
	"1mao/delivery/rest/routes"
	bookingRepository "1mao/internal/booking/repository"
	bookingService "1mao/internal/booking/service"
	clientService "1mao/internal/client/service"
	"1mao/internal/middleware"
	notificationRepository "1mao/internal/notification/repository"
	"1mao/internal/notification/websocket"
	"1mao/internal/payment/repository"
	"1mao/internal/payment/service"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// SetupRoutes configura todas as rotas do sistema
func SetupRoutes(db *gorm.DB, clientService *clientService.ClientService) *mux.Router {
	
	router := mux.NewRouter()

	// Criar repositório de mensagens
	messageRepo := notificationRepository.NewMessageRepository(db)
	bookingService := bookingService.NewBookingService(bookingRepository.NewBookingRepository(db))

	paymentService := service.NewPaymentService(repository.NewPaymentRepository(db), os.Getenv("STRIPE_KEY"))
	// Criar Hub com repositório de mensagens
	hub := websocket.NewHub(messageRepo)
	go hub.Run()

	// Middlewares globais
	router.Use(middleware.LoggerMiddleware)
	router.Use(middleware.RateLimitMiddleware)
	router.Use(middleware.CircuitBreakerMiddleware)

	// Rota do Swagger (DOC)
	routes.SwaggerRouter(router)
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
	// Rotas de agendamento
	routes.BookingRoutes(router, bookingService)
	// Rotas de pagamento
	routes.PaymentRoutes(router, &paymentService )

	return router
}
