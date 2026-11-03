package service

import (
	"context"
	"regexp"
	"strings"

	"github.com/Web-Zefir/secaudit-hub/backend/internal/models"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/repository"
)

type ProjectService struct {
	projects *repository.ProjectRepo
}

func NewProjectService(projects *repository.ProjectRepo) *ProjectService {
	return &ProjectService{projects: projects}
}

type CreateProjectInput struct {
	Name        string
	Description string
	RepoURL     string
	Language    string
}

func (s *ProjectService) Create(ctx context.Context, userID string, in CreateProjectInput) (*models.Project, error) {
	name := strings.TrimSpace(in.Name)
	if len(name) < 2 || len(name) > 80 {
		return nil, ErrInvalidInput
	}
	language := strings.TrimSpace(in.Language)
	if language == "" {
		language = "Mixed"
	}
	slug := slugify(name)
	base := slug
	for i := 2; i < 50; i++ {
		exists, err := s.projects.SlugExists(ctx, userID, slug)
		if err != nil {
			return nil, err
		}
		if !exists {
			break
		}
		slug = base + "-" + itoa(i)
	}
	project := &models.Project{
		UserID:      userID,
		Name:        name,
		Slug:        slug,
		Description: strings.TrimSpace(in.Description),
		RepoURL:     strings.TrimSpace(in.RepoURL),
		Language:    language,
	}
	if err := s.projects.Create(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) List(ctx context.Context, userID string) ([]models.Project, error) {
	items, err := s.projects.ListByUser(ctx, userID)
	if items == nil {
		items = []models.Project{}
	}
	return items, err
}

func (s *ProjectService) Get(ctx context.Context, userID, id string) (*models.Project, error) {
	return s.projects.ByID(ctx, userID, id)
}

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(name string) string {
	s := slugRe.ReplaceAllString(strings.ToLower(strings.TrimSpace(name)), "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "project"
	}
	return s
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [8]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
