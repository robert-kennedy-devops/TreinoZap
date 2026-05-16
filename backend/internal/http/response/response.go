package response

import (
	"encoding/json"
	"net/http"
)

type envelope struct {
	Data any `json:"data,omitempty"`
}

type errorBody struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(envelope{Data: data})
}

func Error(w http.ResponseWriter, status int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorEnvelope{
		Error: errorBody{Message: message, Code: code},
	})
}

func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message, "BAD_REQUEST")
}

func Unauthorized(w http.ResponseWriter) {
	Error(w, http.StatusUnauthorized, "não autenticado", "UNAUTHORIZED")
}

func Forbidden(w http.ResponseWriter) {
	Error(w, http.StatusForbidden, "acesso negado", "FORBIDDEN")
}

func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message, "NOT_FOUND")
}

func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, message, "CONFLICT")
}

func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "erro interno do servidor", "INTERNAL_ERROR")
}
