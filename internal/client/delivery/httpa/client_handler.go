package httpa

import (
	"1mao/internal/client/domain"
	"1mao/internal/client/service"
	"encoding/json"
	"net/http"
)

type ClientHandler struct {
	authService service.ClientService
}

func NewClientHandler(authService service.ClientService) *ClientHandler {
	return &ClientHandler{authService: authService}
}

func (h *ClientHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user domain.Client
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	if err := h.authService.Register(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

}

func (h *ClientHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})

}

func (h *ClientHandler) GetProfile(w http.ResponseWriter, r *http.Request){
	userID, ok := r.Context().Value("userID").(uint)
	if !ok{
		http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *ClientHandler) GetAllUsers(w http.ResponseWriter, r *http.Request){
	users, err := h.authService.GetAllUsers()
	if err != nil {
		http.Error(w, "Erro ao buscar usuários", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *ClientHandler) ForgotPassword(w http.ResponseWriter, r *http.Request){
	var req struct {
		Email string `json:"email"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	message, err := h.authService.ForgotPassword(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}