package routes

import (
	_ "1mao/docs" // Import gerado pelo swag
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
)

func SwaggerRoutes() *mux.Router {
    router := mux.NewRouter()

    // Rota da documentação Swagger
	router.Handle("/swagger/doc.json", http.FileServer(http.Dir("./docs"))).Methods("GET")
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	    
	return router
}