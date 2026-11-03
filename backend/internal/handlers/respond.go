package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Web-Zefir/secaudit-hub/backend/internal/middleware"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/repository"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/service"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<20))
	if err := dec.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	return true
}

func currentUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	session, ok := middleware.Current(r.Context())
	if !ok || session.UserID == "" {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return "", false
	}
	return session.UserID, true
}

func mapServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrEmailTaken):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, repository.ErrNotFound):
		writeError(w, http.StatusNotFound, "not found")
	default:
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}
