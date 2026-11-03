package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"

	"github.com/Web-Zefir/secaudit-hub/backend/internal/config"
	appmw "github.com/Web-Zefir/secaudit-hub/backend/internal/middleware"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/service"
)

type Deps struct {
	Config    config.Config
	DB        *gorm.DB
	Auth      *service.AuthService
	Projects  *service.ProjectService
	Scans     *service.ScanService
	Dashboard *service.DashboardService
}

func NewRouter(deps Deps) http.Handler {
	auth := NewAuthHandler(deps.Auth, deps.Config)
	projects := NewProjectHandler(deps.Projects)
	scans := NewScanHandler(deps.Scans)
	dashboard := NewDashboardHandler(deps.Dashboard, deps.DB)

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(60 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", dashboard.Health)
		r.Get("/modules", dashboard.Modules)
		r.Post("/auth/register", auth.Register)
		r.Post("/auth/login", auth.Login)
		r.Post("/auth/demo", auth.Demo)
		r.Post("/auth/logout", auth.Logout)
		r.Group(func(r chi.Router) {
			r.Use(appmw.Authenticate(deps.Config))
			r.Get("/auth/me", auth.Me)
			r.Get("/dashboard", dashboard.Summary)
			r.Get("/projects", projects.List)
			r.Post("/projects", projects.Create)
			r.Get("/projects/{projectID}", projects.Get)
			r.Get("/projects/{projectID}/scans", scans.List)
			r.Post("/projects/{projectID}/scans", scans.Create)
			r.Get("/scans/{scanID}", scans.Get)
			r.Patch("/findings/{findingID}/status", scans.UpdateFinding)
		})
	})
	return r
}
