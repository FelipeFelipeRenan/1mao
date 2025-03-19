package routes

import (
	"1mao/internal/middleware"
	"1mao/internal/professional/delivery/httpa"
	"1mao/internal/professional/repository"
	"1mao/internal/professional/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// ProfessionalRoutes configura as rotas para profissionais
func ProfessionalRoutes(router *mux.Router, db *gorm.DB) {
	professionalRepo := repository.NewProfessionalRepository(db)
	professionalService := service.NewProfessionalService(professionalRepo)
	professionalHandler := httpa.NewProfessionalHandler(professionalService)

	// Rotas públicas
	router.HandleFunc("/professionals", professionalHandler.GetAllProfessionals).Methods("GET")
	router.HandleFunc("/professional/{id}", professionalHandler.GetProfessionalByID).Methods("GET")
	router.HandleFunc("/professional/register", professionalHandler.Register).Methods("POST")
	router.HandleFunc("/professional/login", professionalHandler.Login).Methods("POST")

	// Rotas protegidas (somente para profissionais autenticados)
	authRouter := router.PathPrefix("/professional").Subrouter()
	authRouter.Use(middleware.AuthMiddleware("professional")) // Middleware agora aceita roles separadas sem precisar de slice
	// Exemplo de rota autenticada (descomentar caso seja necessário)
	// authRouter.HandleFunc("/dashboard", professionalHandler.Dashboard).Methods("GET")
}
