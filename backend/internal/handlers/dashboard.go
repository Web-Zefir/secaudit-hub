package handlers

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/Web-Zefir/secaudit-hub/backend/internal/modules"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/service"
)

type DashboardHandler struct {
	dashboard *service.DashboardService
	db        *gorm.DB
}

func NewDashboardHandler(dashboard *service.DashboardService, db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{dashboard: dashboard, db: db}
}

func (h *DashboardHandler) Health(w http.ResponseWriter, r *http.Request) {
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.PingContext(r.Context()) != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *DashboardHandler) Modules(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, modules.All())
}

func (h *DashboardHandler) Summary(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}
	summary, err := h.dashboard.Build(r.Context(), userID)
	if err != nil {
		mapServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}
