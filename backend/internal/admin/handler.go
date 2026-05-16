package admin

import (
	"net/http"

	"github.com/treinozap/backend/internal/client"
	"github.com/treinozap/backend/internal/http/response"
	"github.com/treinozap/backend/internal/trainer"
	"github.com/treinozap/backend/internal/whatsapp"
)

type Handler struct {
	connMgr              whatsapp.ConnectionManager
	trainerSvc           *trainer.Service
	clientSvc            *client.Service
	whatsAppAdminEnabled bool
}

func NewHandler(connMgr whatsapp.ConnectionManager, trainerSvc *trainer.Service, clientSvc *client.Service, whatsAppAdminEnabled bool) *Handler {
	return &Handler{
		connMgr:              connMgr,
		trainerSvc:           trainerSvc,
		clientSvc:            clientSvc,
		whatsAppAdminEnabled: whatsAppAdminEnabled,
	}
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	if !h.whatsAppAdminEnabled {
		response.Forbidden(w)
		return
	}
	status, err := h.connMgr.Status(r.Context())
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, status)
}

func (h *Handler) QRCode(w http.ResponseWriter, r *http.Request) {
	if !h.whatsAppAdminEnabled {
		response.Forbidden(w)
		return
	}
	qr, err := h.connMgr.QRCode(r.Context())
	if err != nil {
		response.Error(w, http.StatusConflict, err.Error(), "QR_ERROR")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"qr_code": qr})
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	if !h.whatsAppAdminEnabled {
		response.Forbidden(w)
		return
	}
	if err := h.connMgr.Connect(r.Context()); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error(), "CONNECT_ERROR")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "conexão iniciada"})
}

func (h *Handler) Disconnect(w http.ResponseWriter, r *http.Request) {
	if !h.whatsAppAdminEnabled {
		response.Forbidden(w)
		return
	}
	if err := h.connMgr.Disconnect(r.Context()); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error(), "DISCONNECT_ERROR")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "desconectado"})
}

func (h *Handler) ListTrainers(w http.ResponseWriter, r *http.Request) {
	trainers, err := h.trainerSvc.ListAll(r.Context())
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, trainers)
}

func (h *Handler) ListClients(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	clients, err := h.clientSvc.ListAllGlobal(r.Context(), search)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, clients)
}
