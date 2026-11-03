package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Web-Zefir/secaudit-hub/backend/internal/models"
)

type ProjectRepo struct{ db *gorm.DB }

func NewProjectRepo(db *gorm.DB) *ProjectRepo { return &ProjectRepo{db: db} }

func (r *ProjectRepo) Create(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *ProjectRepo) ListByUser(ctx context.Context, userID string) ([]models.Project, error) {
	var items []models.Project
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&items).Error
	return items, err
}

func (r *ProjectRepo) ByID(ctx context.Context, userID, id string) (*models.Project, error) {
	var project models.Project
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &project, err
}

func (r *ProjectRepo) SlugExists(ctx context.Context, userID, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Project{}).Where("user_id = ? AND slug = ?", userID, slug).Count(&count).Error
	return count > 0, err
}

func (r *ProjectRepo) TouchScan(ctx context.Context, id string, score int, at interface{}) error {
	return r.db.WithContext(ctx).Model(&models.Project{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_security_score": score,
		"last_scan_at":        at,
	}).Error
}
