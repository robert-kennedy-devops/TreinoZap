package trainer

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/treinozap/backend/internal/http/middleware"
	"github.com/treinozap/backend/internal/http/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "corpo da requisição inválido")
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		response.BadRequest(w, "nome, e-mail e senha são obrigatórios")
		return
	}
	if len(req.Password) < 6 {
		response.BadRequest(w, "senha deve ter pelo menos 6 caracteres")
		return
	}

	t, err := h.service.Register(r.Context(), req)
	if errors.Is(err, ErrEmailInUse) {
		response.Conflict(w, "e-mail já cadastrado")
		return
	}
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusCreated, t)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "corpo da requisição inválido")
		return
	}

	if req.Email == "" || req.Password == "" {
		response.BadRequest(w, "e-mail e senha são obrigatórios")
		return
	}

	res, err := h.service.Login(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error(), "INVALID_CREDENTIALS")
		return
	}

	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := middleware.TrainerIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	t, err := h.service.GetByID(r.Context(), trainerID)
	if err != nil {
		response.NotFound(w, "treinador não encontrado")
		return
	}

	response.JSON(w, http.StatusOK, t)
}
