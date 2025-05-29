/*package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hashmi846003/P.M.S/internal/models"
	"gorm.io/gorm"
)

// UserRepository interface defines the user-related database operations
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// GetByID retrieves a user by their UUID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email address
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Create inserts a new user into the database
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(user).Error
}
// File: internal/repository/user.go (add these methods)
// Add to existing UserRepository struct

func (r *UserRepository) GetPendingUsers() ([]models.User, error) {
	var users []models.User
	if err := r.db.Where("status = ?", "pending").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) ApproveUser(id uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("status", "approved").Error
}

func (r *UserRepository) RejectUser(id uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("status", "rejected").Error
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}*/
/*package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hashmi846003/P.M.S/internal/models"
	"gorm.io/gorm"
)

// Updated UserRepository interface with all required methods
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	GetPendingUsers(ctx context.Context) ([]models.User, error)
	ApproveUser(ctx context.Context, id uuid.UUID) error
	RejectUser(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, user *models.User) error
	
}

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(user).Error
}

// New methods implemented on the concrete type
/*func (r *userRepository) GetPendingUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Where("status = ?", "pending").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}*/
/*
func (r *userRepository) ApproveUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("status", "approved").Error
}

func (r *userRepository) RejectUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("status", "rejected").Error
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
func (r *userRepository) GetPendingUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Where("status = ?", models.UserStatusPending).
		Find(&users).
		Error
	return users, err
}

func (r *userRepository) ApproveUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("status", models.UserStatusApproved).
		Error
}

func (r *userRepository) RejectUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("status", models.UserStatusRejected).
		Error
}*/
/*
package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hashmi846003/P.M.S/internal/models"
	"gorm.io/gorm"
)

// UserRepository interface defines all required user operations
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	GetPendingUsers(ctx context.Context) ([]models.User, error)
	ApproveUser(ctx context.Context, id uuid.UUID) error
	RejectUser(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, user *models.User) error
}

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetPendingUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Where("status = ?", models.UserStatusPending).
		Find(&users).
		Error
	return users, err
}

func (r *userRepository) ApproveUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("status", models.UserStatusApproved).
		Error
}

func (r *userRepository) RejectUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("status", models.UserStatusRejected).
		Error
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(user).Error
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

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	GetPendingUsers(ctx context.Context) ([]models.User, error)
	ApproveUser(ctx context.Context, id uuid.UUID) error
	RejectUser(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, user *models.User) error
	SetCurrentWorkspace(ctx context.Context, userID, workspaceID uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetPendingUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Where("status = ?", models.UserStatusPending).
		Find(&users).
		Error
	return users, err
}

func (r *userRepository) ApproveUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("status", models.UserStatusApproved).
		Error
}

func (r *userRepository) RejectUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("status", models.UserStatusRejected).
		Error
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) SetCurrentWorkspace(
	ctx context.Context, 
	userID, workspaceID uuid.UUID,
) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("current_workspace_id", workspaceID).
		Error
}