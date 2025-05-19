//package repository
/*
import (
	"context"
	"github.com/google/uuid"
	"github.com/hashmi846003/P.M.S/internal/models"
	"gorm.io/gorm"
)

type DiscussionRepository interface {
	Create(ctx context.Context, discussion *models.Discussion) error
	FindByPageID(ctx context.Context, pageID uuid.UUID) ([]models.Discussion, error)
	GetByPageID(ctx context.Context, pageID uuid.UUID) ([]models.Discussion, error)
}

func (d DiscussionRepository) GetByPageID(context context.Context, pageID uuid.UUID) (any, error) {
	panic("unimplemented")
}

type discussionRepository struct {
	db *gorm.DB
}

func NewDiscussionRepository(db *gorm.DB) DiscussionRepository {
	return &discussionRepository{db: db}
}

func (r *discussionRepository) Create(ctx context.Context, discussion *models.Discussion) error {
	return r.db.WithContext(ctx).Create(discussion).Error
}

func (r *discussionRepository) FindByPageID(ctx context.Context, pageID uuid.UUID) ([]models.Discussion, error) {
	var discussions []models.Discussion
	err := r.db.WithContext(ctx).Where("page_id = ?", pageID).Find(&discussions).Error
	return discussions, err
}
*/
package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/hashmi846003/P.M.S/internal/models"
	"gorm.io/gorm"
)

type DiscussionRepository interface {
	Create(ctx context.Context, discussion *models.Discussion) error
	GetByPageID(ctx context.Context, pageID uuid.UUID) ([]models.Discussion, error)
}

type discussionRepository struct {
	db *gorm.DB
}

func NewDiscussionRepository(db *gorm.DB) DiscussionRepository {
	return &discussionRepository{db: db}
}

func (r *discussionRepository) Create(ctx context.Context, discussion *models.Discussion) error {
	return r.db.WithContext(ctx).Create(discussion).Error
}

func (r *discussionRepository) GetByPageID(ctx context.Context, pageID uuid.UUID) ([]models.Discussion, error) {
	var discussions []models.Discussion
	err := r.db.WithContext(ctx).
		Where("page_id = ?", pageID).
		Order("created_at DESC").
		Find(&discussions).
		Error
	return discussions, err
}