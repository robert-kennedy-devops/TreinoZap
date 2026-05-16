package workout

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/treinozap/backend/internal/client"
	"github.com/treinozap/backend/internal/http/middleware"
	"github.com/treinozap/backend/internal/http/response"
	"github.com/treinozap/backend/internal/message"
	"github.com/treinozap/backend/internal/whatsapp"
)

type Handler struct {
	service     *Service
	clientSvc   *client.Service
	messageRepo *message.Repository
	sender      whatsapp.Sender
}

func NewHandler(service *Service, clientSvc *client.Service, messageRepo *message.Repository, sender whatsapp.Sender) *Handler {
	return &Handler{
		service:     service,
		clientSvc:   clientSvc,
		messageRepo: messageRepo,
		sender:      sender,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := middleware.TrainerIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	clientID, err := uuid.Parse(chi.URLParam(r, "clientId"))
	if err != nil {
		response.BadRequest(w, "clientId inválido")
		return
	}

	// Verify client belongs to trainer
	if _, err := h.clientSvc.GetByID(r.Context(), clientID, trainerID); err != nil {
		response.NotFound(w, "cliente não encontrado")
		return
	}

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "corpo da requisição inválido")
		return
	}
	if req.Name == "" {
		response.BadRequest(w, "nome é obrigatório")
		return
	}

	wo, err := h.service.Create(r.Context(), trainerID, clientID, req)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusCreated, wo)
}

func (h *Handler) ListByClient(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := middleware.TrainerIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	clientID, err := uuid.Parse(chi.URLParam(r, "clientId"))
	if err != nil {
		response.BadRequest(w, "clientId inválido")
		return
	}

	workouts, err := h.service.ListByClient(r.Context(), clientID, trainerID)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, workouts)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := middleware.TrainerIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "id inválido")
		return
	}

	wo, err := h.service.GetByID(r.Context(), id, trainerID)
	if errors.Is(err, ErrNotFound) {
		response.NotFound(w, "treino não encontrado")
		return
	}
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, wo)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := middleware.TrainerIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "id inválido")
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "corpo da requisição inválido")
		return
	}

	wo, err := h.service.Update(r.Context(), id, trainerID, req)
	if errors.Is(err, ErrNotFound) {
		response.NotFound(w, "treino não encontrado")
		return
	}
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, wo)
}

func (h *Handler) Archive(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := middleware.TrainerIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "id inválido")
		return
	}

	if err := h.service.Archive(r.Context(), id, trainerID); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, "treino não encontrado")
			return
		}
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "treino arquivado"})
}

func (h *Handler) Activate(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := middleware.TrainerIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "id inválido")
		return
	}

	wo, err := h.service.Activate(r.Context(), id, trainerID)
	if errors.Is(err, ErrNotFound) {
		response.NotFound(w, "treino não encontrado")
		return
	}
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, wo)
}

func (h *Handler) SendWhatsApp(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := middleware.TrainerIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "id inválido")
		return
	}

	wo, err := h.service.GetByID(r.Context(), id, trainerID)
	if errors.Is(err, ErrNotFound) {
		response.NotFound(w, "treino não encontrado")
		return
	}
	if err != nil {
		response.InternalError(w)
		return
	}

	c, err := h.clientSvc.GetByID(r.Context(), wo.ClientID, trainerID)
	if err != nil {
		response.NotFound(w, "cliente não encontrado")
		return
	}

	text := FormatWorkoutForWhatsApp(c.Name, wo)

	if err := h.sender.SendText(r.Context(), c.Phone, text); err != nil {
		log.Printf("[workout] erro ao enviar WhatsApp para %s: %v", c.Phone, err)
		response.InternalError(w)
		return
	}

	clientID := c.ID
	msg := message.NewOutbound(&trainerID, &clientID, c.Phone, text, "treino_manual")
	if err := h.messageRepo.Save(r.Context(), msg); err != nil {
		log.Printf("[workout] erro ao salvar mensagem: %v", err)
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "treino enviado com sucesso"})
}
