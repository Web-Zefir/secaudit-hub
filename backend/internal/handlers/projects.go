package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Web-Zefir/secaudit-hub/backend/internal/service"
)

type ProjectHandler struct {
	projects *service.ProjectService
}

func NewProjectHandler(projects *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projects: projects}
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}
	items, err := h.projects.List(r.Context(), userID)
	if err != nil {
		mapServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		RepoURL     string `json:"repoUrl"`
		Language    string `json:"language"`
	}
	if !decode(w, r, &body) {
		return
	}
	project, err := h.projects.Create(r.Context(), userID, service.CreateProjectInput{
		Name: body.Name, Description: body.Description, RepoURL: body.RepoURL, Language: body.Language,
	})
	if err != nil {
		mapServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}
	project, err := h.projects.Get(r.Context(), userID, chi.URLParam(r, "projectID"))
	if err != nil {
		mapServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, project)
}
