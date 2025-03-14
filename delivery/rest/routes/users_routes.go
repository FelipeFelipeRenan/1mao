package routes

import (
	"1mao/internal/middleware"
	"1mao/internal/user/delivery/httpa"
	"1mao/internal/user/service"

	"github.com/gorilla/mux"
)


func UserRoutes(r *mux.Router, userService *service.AuthService){

	userHandler := httpa.NewUserHandler(*userService)

	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")

	authRouter := r.PathPrefix("/client").Subrouter()
	authRouter.Use(middleware.AuthMiddleware([]string{"client"}))
	authRouter.HandleFunc("/me", userHandler.GetProfile).Methods("GET")
}