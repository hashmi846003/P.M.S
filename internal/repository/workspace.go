/*
package repository

import (
	"github.com/google/uuid"
	"github.com/hashmi846003/P.M.S/internal/models"
	"gorm.io/gorm"
)


type WorkspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

func (r *WorkspaceRepository) Create(workspace *models.Workspace) error {
	return r.db.Create(workspace).Error
}

func (r *WorkspaceRepository) GetByID(id uuid.UUID) (*models.Workspace, error) {
	var workspace models.Workspace
	if err := r.db.First(&workspace, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (r *WorkspaceRepository) GetByOwner(ownerID uuid.UUID) ([]models.Workspace, error) {
	var workspaces []models.Workspace
	if err := r.db.Find(&workspaces, "owner_id = ?", ownerID).Error; err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (r *WorkspaceRepository) Update(workspace *models.Workspace) error {
	return r.db.Save(workspace).Error
}

func (r *WorkspaceRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Workspace{}, "id = ?", id).Error
}*/
package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hashmi846003/P.M.S/internal/models"
	"gorm.io/gorm"
)

type WorkspaceRepository interface {
	Create(ctx context.Context, workspace *models.Workspace) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Workspace, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Workspace, error)
	Update(ctx context.Context, workspace *models.Workspace) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddUserToWorkspace(ctx context.Context, workspaceID, userID uuid.UUID, isOwner bool) error
	IsWorkspaceOwner(ctx context.Context, workspaceID, userID uuid.UUID) (bool, error)
	UserHasAccess(ctx context.Context, workspaceID, userID uuid.UUID) (bool, error)
}

type workspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) WorkspaceRepository {
	return &workspaceRepository{db: db}
}

func (r *workspaceRepository) Create(ctx context.Context, workspace *models.Workspace) error {
	workspace.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(workspace).Error
}

func (r *workspaceRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Workspace, error) {
	var workspace models.Workspace
	err := r.db.WithContext(ctx).First(&workspace, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &workspace, nil
}

func (r *workspaceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Workspace, error) {
	var workspaces []*models.Workspace
	err := r.db.WithContext(ctx).
		Joins("JOIN user_workspaces ON user_workspaces.workspace_id = workspaces.id").
		Where("user_workspaces.user_id = ?", userID).
		Find(&workspaces).
		Error
	return workspaces, err
}

func (r *workspaceRepository) Update(ctx context.Context, workspace *models.Workspace) error {
	workspace.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(workspace).Error
}

func (r *workspaceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Workspace{}, "id = ?", id).Error
}

func (r *workspaceRepository) AddUserToWorkspace(
	ctx context.Context, 
	workspaceID, userID uuid.UUID, 
	isOwner bool,
) error {
	userWorkspace := &models.UserWorkspace{
		UserID:      userID,
		WorkspaceID: workspaceID,
		IsOwner:     isOwner,
		CreatedAt:   time.Now(),
	}
	return r.db.WithContext(ctx).Create(userWorkspace).Error
}

func (r *workspaceRepository) IsWorkspaceOwner(
	ctx context.Context, 
	workspaceID, userID uuid.UUID,
) (bool, error) {
	var userWorkspace models.UserWorkspace
	err := r.db.WithContext(ctx).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		First(&userWorkspace).
		Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return userWorkspace.IsOwner, nil
}

func (r *workspaceRepository) UserHasAccess(
	ctx context.Context, 
	workspaceID, userID uuid.UUID,
) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.UserWorkspace{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Count(&count).
		Error
	
	return count > 0, err
}