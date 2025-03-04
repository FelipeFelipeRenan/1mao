package httpa

import (
	"1mao/internal/user/domain"
	"1mao/internal/user/service"
	"encoding/json"
	"net/http"
)


type UserHandler struct {
	authService service.AuthService
}

func NewUserHandler(authService service.AuthService) *UserHandler{
	return &UserHandler{authService: authService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request){
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil{
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	if err := h.authService.Register(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message":"Usuário registrado com sucesso"})

}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request){
	var creds struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil{
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token":token})


}