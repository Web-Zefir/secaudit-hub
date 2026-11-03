package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Web-Zefir/secaudit-hub/backend/internal/models"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/service"
)

type ScanHandler struct {
	scans *service.ScanService
}

func NewScanHandler(scans *service.ScanService) *ScanHandler {
	return &ScanHandler{scans: scans}
}

func (h *ScanHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}
	items, err := h.scans.List(r.Context(), userID, chi.URLParam(r, "projectID"))
	if err != nil {
		mapServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *ScanHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}
	var body struct {
		TargetSource string            `json:"targetSource"`
		RepoURL      string            `json:"repoUrl"`
		Files        []models.CodeFile `json:"files"`
	}
	if !decode(w, r, &body) {
		return
	}
	job, err := h.scans.Start(r.Context(), userID, chi.URLParam(r, "projectID"), service.StartScanInput{
		TargetSource: body.TargetSource,
		RepoURL:      body.RepoURL,
		Files:        body.Files,
	})
	if err != nil {
		mapServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, job)
}

func (h *ScanHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}
	job, findings, err := h.scans.Get(r.Context(), userID, chi.URLParam(r, "scanID"))
	if err != nil {
		mapServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"scan": job, "findings": findings})
}

func (h *ScanHandler) UpdateFinding(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if !decode(w, r, &body) {
		return
	}
	item, err := h.scans.UpdateFindingStatus(r.Context(), userID, chi.URLParam(r, "findingID"), body.Status)
	if err != nil {
		mapServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
