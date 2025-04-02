package routes

import (
	"1mao/internal/middleware"
	"1mao/internal/client/delivery/httpa"
	"1mao/internal/client/service"

	"github.com/gorilla/mux"
)

// UserRoutes configura as rotas para usuários
func UserRoutes(r *mux.Router, userService *service.ClientService) {
	userHandler := httpa.NewClientHandler(*userService)

	// Rotas públicas
	r.HandleFunc("/client/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/client/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/client/users", userHandler.GetAllUsers).Methods("GET")

	// Rotas protegidas (somente para clientes autenticados)
	authRouter := r.PathPrefix("/client").Subrouter()
	authRouter.Use(middleware.AuthMiddleware("client"))
	authRouter.HandleFunc("/me", userHandler.GetProfile).Methods("GET")
}
