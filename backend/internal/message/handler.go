package message

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/treinozap/backend/internal/client"
	"github.com/treinozap/backend/internal/http/middleware"
	"github.com/treinozap/backend/internal/http/response"
)

type Handler struct {
	repo      *Repository
	clientSvc *client.Service
}

func NewHandler(repo *Repository, clientSvc *client.Service) *Handler {
	return &Handler{repo: repo, clientSvc: clientSvc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := middleware.TrainerIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	direction := r.URL.Query().Get("direction")

	messages, total, err := h.repo.List(r.Context(), ListFilter{
		TrainerID: &trainerID,
		Direction: direction,
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		response.InternalError(w)
		return
	}
	if messages == nil {
		messages = []Message{}
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data":        messages,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + pageSize - 1) / pageSize,
	})
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

	// Verify the client belongs to the authenticated trainer before listing messages.
	if _, err := h.clientSvc.GetByID(r.Context(), clientID, trainerID); err != nil {
		if errors.Is(err, client.ErrNotFound) {
			response.NotFound(w, "cliente não encontrado")
		} else {
			response.InternalError(w)
		}
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	messages, total, err := h.repo.List(r.Context(), ListFilter{
		TrainerID: &trainerID,
		ClientID:  &clientID,
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		response.InternalError(w)
		return
	}
	if messages == nil {
		messages = []Message{}
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data":        messages,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + pageSize - 1) / pageSize,
	})
}
