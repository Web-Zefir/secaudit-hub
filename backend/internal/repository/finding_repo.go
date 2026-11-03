package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Web-Zefir/secaudit-hub/backend/internal/models"
)

type FindingRepo struct{ db *gorm.DB }

func NewFindingRepo(db *gorm.DB) *FindingRepo { return &FindingRepo{db: db} }

func (r *FindingRepo) CreateBatch(ctx context.Context, items []models.Finding) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(items, 100).Error
}

func (r *FindingRepo) ByScan(ctx context.Context, userID, scanID string) ([]models.Finding, error) {
	var items []models.Finding
	err := r.db.WithContext(ctx).
		Select("findings.*").
		Joins("JOIN projects ON projects.id = findings.project_id").
		Where("findings.scan_job_id = ? AND projects.user_id = ?", scanID, userID).
		Order("CASE findings.severity WHEN 'CRITICAL' THEN 1 WHEN 'HIGH' THEN 2 WHEN 'MEDIUM' THEN 3 ELSE 4 END, findings.file_path, findings.line_number").
		Find(&items).Error
	return items, err
}

func (r *FindingRepo) UpdateStatus(ctx context.Context, userID, id, status string) (*models.Finding, error) {
	var item models.Finding
	err := r.db.WithContext(ctx).
		Select("findings.*").
		Joins("JOIN projects ON projects.id = findings.project_id").
		Where("findings.id = ? AND projects.user_id = ?", id, userID).
		First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	item.Status = status
	if err := r.db.WithContext(ctx).Model(&models.Finding{}).Where("id = ?", item.ID).Update("status", status).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

type SeverityCount struct {
	Severity string
	Count    int64
}

type CategoryCount struct {
	Category string
	Count    int64
}

func (r *FindingRepo) OpenSeverityCounts(ctx context.Context, userID string) ([]SeverityCount, error) {
	var rows []SeverityCount
	err := r.db.WithContext(ctx).Model(&models.Finding{}).
		Select("findings.severity as severity, count(*) as count").
		Joins("JOIN projects ON projects.id = findings.project_id").
		Where("projects.user_id = ? AND findings.status = ?", userID, models.FindingOpen).
		Group("findings.severity").
		Scan(&rows).Error
	return rows, err
}

func (r *FindingRepo) OpenCategoryCounts(ctx context.Context, userID string) ([]CategoryCount, error) {
	var rows []CategoryCount
	err := r.db.WithContext(ctx).Model(&models.Finding{}).
		Select("findings.category as category, count(*) as count").
		Joins("JOIN projects ON projects.id = findings.project_id").
		Where("projects.user_id = ? AND findings.status = ?", userID, models.FindingOpen).
		Group("findings.category").
		Order("count desc").
		Scan(&rows).Error
	return rows, err
}
