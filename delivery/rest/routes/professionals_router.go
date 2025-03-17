package routes

import (
	"1mao/internal/professional/delivery/httpa"
	"1mao/internal/professional/repository"
	"1mao/internal/professional/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func ProfessionalRoutes(router *mux.Router, db *gorm.DB){

	professionalRepo := repository.NewProfessionalRepository(db)
	professinonalService := service.NewProfessionalService(professionalRepo)
	professionalHandler := httpa.NewProfessionalHandler(professinonalService)

	router.HandleFunc("/professionals", professionalHandler.GetAllProfessionals).Methods("GET")
	router.HandleFunc("/professional/{id}", professionalHandler.GetProfessionalByID).Methods("GET")
	router.HandleFunc("professional/register", professionalHandler.Register).Methods("POST")


}