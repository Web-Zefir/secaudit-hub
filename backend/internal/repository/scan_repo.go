package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/Web-Zefir/secaudit-hub/backend/internal/models"
)

type ScanRepo struct{ db *gorm.DB }

func NewScanRepo(db *gorm.DB) *ScanRepo { return &ScanRepo{db: db} }

func (r *ScanRepo) Create(ctx context.Context, job *models.ScanJob) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *ScanRepo) ListByProject(ctx context.Context, userID, projectID string) ([]models.ScanJob, error) {
	var items []models.ScanJob
	err := r.db.WithContext(ctx).
		Select("scan_jobs.*").
		Joins("JOIN projects ON projects.id = scan_jobs.project_id").
		Where("scan_jobs.project_id = ? AND projects.user_id = ?", projectID, userID).
		Order("scan_jobs.created_at desc").
		Find(&items).Error
	return items, err
}

func (r *ScanRepo) ByID(ctx context.Context, userID, id string) (*models.ScanJob, error) {
	var job models.ScanJob
	err := r.db.WithContext(ctx).
		Select("scan_jobs.*").
		Joins("JOIN projects ON projects.id = scan_jobs.project_id").
		Where("scan_jobs.id = ? AND projects.user_id = ?", id, userID).
		First(&job).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &job, err
}

func (r *ScanRepo) MarkRunning(ctx context.Context, id string, started time.Time) error {
	return r.db.WithContext(ctx).Model(&models.ScanJob{}).Where("id = ?", id).Updates(map[string]any{
		"status":     models.ScanRunning,
		"started_at": started,
	}).Error
}

func (r *ScanRepo) Finish(ctx context.Context, job *models.ScanJob) error {
	return r.db.WithContext(ctx).Model(&models.ScanJob{}).Where("id = ?", job.ID).Updates(map[string]any{
		"status":         job.Status,
		"security_score": job.SecurityScore,
		"files_scanned":  job.FilesScanned,
		"duration_ms":    job.DurationMs,
		"critical_count": job.CriticalCount,
		"high_count":     job.HighCount,
		"medium_count":   job.MediumCount,
		"low_count":      job.LowCount,
		"error_message":  job.ErrorMessage,
		"finished_at":    job.FinishedAt,
	}).Error
}
