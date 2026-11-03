package service

import (
	"context"

	"github.com/Web-Zefir/secaudit-hub/backend/internal/models"
	"github.com/Web-Zefir/secaudit-hub/backend/internal/repository"
)

type DashboardService struct {
	projects *repository.ProjectRepo
	findings *repository.FindingRepo
}

func NewDashboardService(projects *repository.ProjectRepo, findings *repository.FindingRepo) *DashboardService {
	return &DashboardService{projects: projects, findings: findings}
}

type CategoryStat struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

type Dashboard struct {
	AverageScore *float64         `json:"averageScore"`
	ProjectCount int              `json:"projectCount"`
	OpenCritical int64            `json:"openCritical"`
	OpenHigh     int64            `json:"openHigh"`
	Categories   []CategoryStat   `json:"categories"`
	Projects     []models.Project `json:"projects"`
}

func (s *DashboardService) Build(ctx context.Context, userID string) (*Dashboard, error) {
	projects, err := s.projects.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if projects == nil {
		projects = []models.Project{}
	}
	var sum int
	var scored int
	for _, p := range projects {
		if p.LastSecurityScore != nil {
			sum += *p.LastSecurityScore
			scored++
		}
	}
	var avg *float64
	if scored > 0 {
		value := float64(sum) / float64(scored)
		avg = &value
	}
	severities, err := s.findings.OpenSeverityCounts(ctx, userID)
	if err != nil {
		return nil, err
	}
	var critical, high int64
	for _, row := range severities {
		switch row.Severity {
		case models.SeverityCritical:
			critical = row.Count
		case models.SeverityHigh:
			high = row.Count
		}
	}
	cats, err := s.findings.OpenCategoryCounts(ctx, userID)
	if err != nil {
		return nil, err
	}
	stats := make([]CategoryStat, 0, len(cats))
	for _, c := range cats {
		stats = append(stats, CategoryStat{Category: c.Category, Count: c.Count})
	}
	return &Dashboard{
		AverageScore: avg,
		ProjectCount: len(projects),
		OpenCritical: critical,
		OpenHigh:     high,
		Categories:   stats,
		Projects:     projects,
	}, nil
}
