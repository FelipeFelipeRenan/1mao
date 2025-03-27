package httpa

import (
	"1mao/internal/client/domain"
	"1mao/internal/client/service"
	"encoding/json"
	"net/http"
)

// RegisterRequest define a estrutura para registro de clientes
//	@Description	Dados necessários para registrar um novo cliente
type RegisterRequest struct {
	Name     string `json:"name" example:"João Silva"`
	Email    string `json:"email" example:"cliente@example.com"`
	Password string `json:"password" example:"senhaSegura123"`
	Phone    string `json:"phone" example:"+5511999999999"`
}

// LoginRequest define a estrutura para login de clientes
//	@Description	Credenciais para autenticação do cliente
type LoginRequest struct {
	Email    string `json:"email" example:"cliente@example.com"`
	Password string `json:"password" example:"senhaSegura123"`
}

// LoginResponse define a estrutura da resposta do login
//	@Description	Retorno do endpoint de login contendo o token JWT
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ClientHandler lida com requisições de clientes
//	@title			Client API
//	@version		1.0
//	@description	API para gestão de clientes
type ClientHandler struct {
	authService service.ClientService
}

func NewClientHandler(authService service.ClientService) *ClientHandler {
	return &ClientHandler{authService: authService}
}


// Register godoc
//	@Summary		Registrar novo cliente
//	@Description	Cria uma nova conta de cliente
//	@Tags			Clients
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterRequest	true	"Dados do cliente"
//	@Success		201		{object}	domain.Client
//	@Failure		400		{object}	map[string]string	"Dados inválidos"
//	@Failure		500		{object}	map[string]string	"Erro interno"
//	@Router			/client/register [post]
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

// Login godoc
//	@Summary		Login de cliente
//	@Description	Autentica um cliente e retorna token JWT
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest	true	"Credenciais de login"
//	@Success		200		{object}	LoginResponse
//	@Failure		400		{object}	map[string]string	"Requisição inválida"
//	@Failure		401		{object}	map[string]string	"Credenciais inválidas"
//	@Router			/client/login [post]
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

// GetProfile godoc
//	@Summary		Obter perfil do cliente
//	@Description	Retorna os dados do cliente autenticado
//	@Tags			Clients
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	domain.Client
//	@Failure		401	{object}	map[string]string	"Não autorizado"
//	@Router			/client/me [get]
func (h *ClientHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
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

func (h *ClientHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.authService.GetAllUsers()
	if err != nil {
		http.Error(w, "Erro ao buscar usuários", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *ClientHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
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
