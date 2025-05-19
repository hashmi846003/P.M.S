package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashmi846003/P.M.S/internal/models"
	"gorm.io/gorm"
)

type PageRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID, includeDeleted bool) ([]models.Page, error)
	Create(ctx context.Context, page *models.Page) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Page, error)
	Update(ctx context.Context, page *models.Page) error
	SoftDelete(ctx context.Context, pageID, userID uuid.UUID) error
	Restore(ctx context.Context, pageID, userID uuid.UUID) error
	ToggleFavorite(ctx context.Context, pageID, userID uuid.UUID) (bool, error)
	Duplicate(ctx context.Context, pageID, userID uuid.UUID) (*models.Page, error)
	GetVersions(ctx context.Context, pageID uuid.UUID) ([]models.PageVersion, error)
	CreateVersion(ctx context.Context, version *models.PageVersion) error
	ApplyFormatting(ctx context.Context, pageID, userID uuid.UUID, format string, start, end int) (*models.Page, error)
	SetAlignment(ctx context.Context, pageID, userID uuid.UUID, alignment string) (*models.Page, error)
	AddEmoji(ctx context.Context, pageID, userID uuid.UUID, emoji string) (*models.Page, error)
	CreateShareLink(ctx context.Context, pageID uuid.UUID, token string, expiresAt time.Time) error
}

type pageRepository struct {
	db *gorm.DB
}

func NewPageRepository(db *gorm.DB) PageRepository {
	return &pageRepository{db: db}
}

func (r *pageRepository) FindByUserID(ctx context.Context, userID uuid.UUID, includeDeleted bool) ([]models.Page, error) {
	var pages []models.Page
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if !includeDeleted {
		query = query.Where("is_deleted = false")
	}
	err := query.Find(&pages).Error
	return pages, err
}

func (r *pageRepository) Create(ctx context.Context, page *models.Page) error {
	return r.db.WithContext(ctx).Create(page).Error
}

func (r *pageRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Page, error) {
	var page models.Page
	err := r.db.WithContext(ctx).First(&page, "id = ?", id).Error
	return &page, err
}

func (r *pageRepository) Update(ctx context.Context, page *models.Page) error {
	return r.db.WithContext(ctx).Save(page).Error
}

func (r *pageRepository) SoftDelete(ctx context.Context, pageID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Page{}).
		Where("id = ? AND user_id = ?", pageID, userID).
		Update("is_deleted", true).Error
}

func (r *pageRepository) Restore(ctx context.Context, pageID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Page{}).
		Where("id = ? AND user_id = ?", pageID, userID).
		Update("is_deleted", false).Error
}

func (r *pageRepository) ToggleFavorite(ctx context.Context, pageID, userID uuid.UUID) (bool, error) {
	var page models.Page
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&page, "id = ? AND user_id = ?", pageID, userID).Error; err != nil {
			return err
		}
		return tx.Model(&page).Update("is_favorite", !page.IsFavorite).Error
	})
	return !page.IsFavorite, err
}

func (r *pageRepository) Duplicate(ctx context.Context, pageID, userID uuid.UUID) (*models.Page, error) {
	var newPage *models.Page
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		original, err := r.FindByID(ctx, pageID)
		if err != nil {
			return err
		}

		newPage = &models.Page{
			ID:         uuid.New(),
			Title:      fmt.Sprintf("Copy of %s", original.Title),
			Content:    original.Content,
			ParentID:   original.ParentID,
			UserID:     userID,
			IsDeleted:  false,
			IsFavorite: false,
			Emoji:      original.Emoji,
			Icon:       original.Icon,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		return tx.Create(newPage).Error
	})
	return newPage, err
}

func (r *pageRepository) GetVersions(ctx context.Context, pageID uuid.UUID) ([]models.PageVersion, error) {
	var versions []models.PageVersion
	err := r.db.WithContext(ctx).Where("page_id = ?", pageID).Find(&versions).Error
	return versions, err
}

func (r *pageRepository) CreateVersion(ctx context.Context, version *models.PageVersion) error {
	return r.db.WithContext(ctx).Create(version).Error
}

func (r *pageRepository) ApplyFormatting(ctx context.Context, pageID, userID uuid.UUID, format string, start, end int) (*models.Page, error) {
	var page models.Page
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&page, "id = ? AND user_id = ?", pageID, userID).Error; err != nil {
			return err
		}

		content := []rune(page.Content)
		if start < 0 || end > len(content) || start > end {
			return fmt.Errorf("invalid formatting range")
		}

		var formatted string
		switch format {
		case "bold":
			formatted = fmt.Sprintf("**%s**", string(content[start:end]))
		case "italic":
			formatted = fmt.Sprintf("*%s*", string(content[start:end]))
		case "underline":
			formatted = fmt.Sprintf("<u>%s</u>", string(content[start:end]))
		default:
			return fmt.Errorf("unsupported format type")
		}

		page.Content = string(content[:start]) + formatted + string(content[end:])
		return tx.Save(&page).Error
	})
	return &page, err
}

func (r *pageRepository) SetAlignment(ctx context.Context, pageID, userID uuid.UUID, alignment string) (*models.Page, error) {
	var page models.Page
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&page, "id = ? AND user_id = ?", pageID, userID).Error; err != nil {
			return err
		}

		page.Content = fmt.Sprintf(`<div style="text-align: %s">%s</div>`, alignment, page.Content)
		return tx.Save(&page).Error
	})
	return &page, err
}

func (r *pageRepository) AddEmoji(ctx context.Context, pageID, userID uuid.UUID, emoji string) (*models.Page, error) {
	var page models.Page
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&page, "id = ? AND user_id = ?", pageID, userID).Error; err != nil {
			return err
		}
		page.Emoji = emoji
		return tx.Save(&page).Error
	})
	return &page, err
}

func (r *pageRepository) CreateShareLink(ctx context.Context, pageID uuid.UUID, token string, expiresAt time.Time) error {
	return r.db.WithContext(ctx).Create(&models.ShareLink{
		Token:     token,
		PageID:    pageID,
		ExpiresAt: expiresAt,
	}).Error
}