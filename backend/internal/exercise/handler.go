package exercise

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/treinozap/backend/internal/http/middleware"
	"github.com/treinozap/backend/internal/http/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := middleware.TrainerIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
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

	e, err := h.service.Create(r.Context(), trainerID, req)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusCreated, e)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := middleware.TrainerIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	exercises, total, err := h.service.List(r.Context(), ListFilter{
		TrainerID: trainerID,
		Search:    r.URL.Query().Get("search"),
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		response.InternalError(w)
		return
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data":        exercises,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + pageSize - 1) / pageSize,
	})
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

	e, err := h.service.GetByID(r.Context(), id, trainerID)
	if errors.Is(err, ErrNotFound) {
		response.NotFound(w, "exercício não encontrado")
		return
	}
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, e)
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
	if req.Name == "" {
		response.BadRequest(w, "nome é obrigatório")
		return
	}

	e, err := h.service.Update(r.Context(), id, trainerID, req)
	if errors.Is(err, ErrNotFound) {
		response.NotFound(w, "exercício não encontrado")
		return
	}
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, e)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := h.service.Delete(r.Context(), id, trainerID); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, "exercício não encontrado")
			return
		}
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "exercício removido"})
}
