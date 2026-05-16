package client

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
	if req.Name == "" || req.Phone == "" {
		response.BadRequest(w, "nome e telefone são obrigatórios")
		return
	}

	c, err := h.service.Create(r.Context(), trainerID, req)
	if errors.Is(err, ErrPhoneInUse) {
		response.Conflict(w, "telefone já cadastrado")
		return
	}
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusCreated, c)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := middleware.TrainerIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	search := r.URL.Query().Get("search")

	clients, total, err := h.service.List(r.Context(), ListFilter{
		TrainerID: trainerID,
		Search:    search,
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
	totalPages := (total + pageSize - 1) / pageSize

	response.JSON(w, http.StatusOK, map[string]any{
		"data":        clients,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
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

	c, err := h.service.GetByID(r.Context(), id, trainerID)
	if errors.Is(err, ErrNotFound) {
		response.NotFound(w, "cliente não encontrado")
		return
	}
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, c)
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
	if req.Name == "" || req.Phone == "" {
		response.BadRequest(w, "nome e telefone são obrigatórios")
		return
	}

	c, err := h.service.Update(r.Context(), id, trainerID, req)
	if errors.Is(err, ErrNotFound) {
		response.NotFound(w, "cliente não encontrado")
		return
	}
	if errors.Is(err, ErrPhoneInUse) {
		response.Conflict(w, "telefone já cadastrado")
		return
	}
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, c)
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
			response.NotFound(w, "cliente não encontrado")
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "cliente inativado com sucesso"})
}
