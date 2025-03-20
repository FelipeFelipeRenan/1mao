package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)



var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool{return true},
}

func HandleConnections(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		fmt.Println("Erro ao fazer upgrade para WebSocket: ", err)
		return 
	}
	defer conn.Close()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Erro ao ler mensagem: ", err)
			break
		}
		fmt.Printf("Mensagem recebida: %s", msg)
		conn.WriteMessage(messageType, []byte("Mensagem recebida!"))
	}
}