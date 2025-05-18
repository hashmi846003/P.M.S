package repository

import (
	"context"
	
	"gorm.io/gorm"
	"your-app/internal/models"
)

type PageRepository interface {
	CreatePage(ctx context.Context, page *models.Page) (*models.Page, error)
	GetAllPages(ctx context.Context, userID uint) ([]models.Page, error)
	GetPageByID(ctx context.Context, pageID uint) (*models.Page, error)
	UpdatePage(ctx context.Context, page *models.Page) (*models.Page, error)
	DeletePage(ctx context.Context, pageID uint) error
}

type pageRepository struct {
	db *gorm.DB
}

func NewPageRepository(db *gorm.DB) PageRepository {
	return &pageRepository{db: db}
}

func (r *pageRepository) CreatePage(ctx context.Context, page *models.Page) (*models.Page, error) {
	err := r.db.WithContext(ctx).Create(page).Error
	return page, err
}

func (r *pageRepository) GetAllPages(ctx context.Context, userID uint) ([]models.Page, error) {
	var pages []models.Page
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_deleted = ?", userID, false).Find(&pages).Error
	return pages, err
}

func (r *pageRepository) GetPageByID(ctx context.Context, pageID uint) (*models.Page, error) {
	var page models.Page
	err := r.db.WithContext(ctx).First(&page, pageID).Error
	return &page, err
}

func (r *pageRepository) UpdatePage(ctx context.Context, page *models.Page) (*models.Page, error) {
	err := r.db.WithContext(ctx).Save(page).Error
	return page, err
}

func (r *pageRepository) DeletePage(ctx context.Context, pageID uint) error {
	return r.db.WithContext(ctx).Model(&models.Page{}).
		Where("id = ?", pageID).
		Update("is_deleted", true).Error
}