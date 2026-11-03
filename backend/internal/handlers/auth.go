package handlers

import (
	"net/http"

	"github.com/Web-Zefir/secaudit-hub/backend/internal/config"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/middleware"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
	cfg  config.Config
}

func NewAuthHandler(auth *service.AuthService, cfg config.Config) *AuthHandler {
	return &AuthHandler{auth: auth, cfg: cfg}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !decode(w, r, &body) {
		return
	}
	user, err := h.auth.Register(r.Context(), body.Name, body.Email, body.Password)
	if err != nil {
		mapServiceError(w, err)
		return
	}
	if !h.startSession(w, user.ID, user.Email) {
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !decode(w, r, &body) {
		return
	}
	user, err := h.auth.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		mapServiceError(w, err)
		return
	}
	if !h.startSession(w, user.ID, user.Email) {
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) Demo(w http.ResponseWriter, r *http.Request) {
	user, err := h.auth.Demo(r.Context())
	if err != nil {
		mapServiceError(w, err)
		return
	}
	if !h.startSession(w, user.ID, user.Email) {
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	middleware.ClearSessionCookie(w, h.cfg)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}
	user, err := h.auth.Me(r.Context(), userID)
	if err != nil {
		mapServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) startSession(w http.ResponseWriter, userID, email string) bool {
	token, err := middleware.IssueToken(h.cfg.JWTSecret, userID, email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return false
	}
	middleware.SetSessionCookie(w, h.cfg, token)
	return true
}
