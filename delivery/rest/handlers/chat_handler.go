package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	client "1mao/internal/client/repository"
	"1mao/internal/notification/websocket"
	professional "1mao/internal/professional/repository"

	"github.com/gorilla/mux"
	ws "github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// Configuração do WebSocket
var chatUpgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleClientChatWebSocket(w http.ResponseWriter, r *http.Request, db *gorm.DB, hub *websocket.Hub){
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if !clientExists(db, userID) {
		http.Error(w, "Client nao encontrado", http.StatusNotFound)
		return
	}

	handleChatWebSocket(w,r,userID, hub)
}

func HandleProfessionalChatWebSocket(w http.ResponseWriter, r *http.Request, db *gorm.DB, hub *websocket.Hub){
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if !professionalExists(db, userID) {
		http.Error(w, "Profissional nao encontrado", http.StatusNotFound)
		return
	}

	handleChatWebSocket(w,r,userID, hub)
}


// Manipulador de WebSocket para Chat
func handleChatWebSocket(w http.ResponseWriter, r *http.Request,userID int, hub *websocket.Hub) {

	conn, err := chatUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("❌ Erro ao conectar WebSocket:", err)
		http.Error(w, "Erro ao conectar WebSocket", http.StatusInternalServerError)
		return
	}

	// Criar um novo cliente e registrá-lo no Hub
	client := websocket.NewClient(userID, conn, hub)
	hub.Register <- client

	fmt.Printf("🔗 Cliente %d conectado ao WebSocket do Chat\n", userID)

	// Inicia as goroutines para leitura e escrita
	go client.Listen()
	go client.Write()
}

func clientExists(db *gorm.DB, userID int) bool{
	clientRepo := client.NewUserRepository(db)
	_, err := clientRepo.FindByID(uint(userID))
	return err == nil
}

func professionalExists(db *gorm.DB, userID int) bool{
	clientRepo := professional.NewProfessionalRepository(db)
	_, err := clientRepo.FindByID(uint(userID))
	return err == nil
}