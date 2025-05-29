
/*package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hashmi846003/P.M.S/internal/models"
	"github.com/hashmi846003/P.M.S/internal/repository"
)

type WorkspaceHandler struct {
	workspaceRepo *repository.WorkspaceRepository
	userRepo      *repository.UserRepository
}

func NewWorkspaceHandler(workspaceRepo *repository.WorkspaceRepository, userRepo *repository.UserRepository) *WorkspaceHandler {
	return &WorkspaceHandler{workspaceRepo: workspaceRepo, userRepo: userRepo}
}

func (h *WorkspaceHandler) CreateWorkspace(c *gin.Context) {
	var workspace models.Workspace
	if err := c.ShouldBindJSON(&workspace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Get current user from context (set by auth middleware)
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user := currentUser.(*models.User)

	workspace.OwnerID = user.ID

	if err := h.workspaceRepo.Create(&workspace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workspace"})
		return
	}

	c.JSON(http.StatusCreated, workspace)
}

func (h *WorkspaceHandler) ListWorkspaces(c *gin.Context) {
	// Get current user from context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user := currentUser.(*models.User)

	workspaces, err := h.workspaceRepo.GetByOwner(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch workspaces"})
		return
	}
	c.JSON(http.StatusOK, workspaces)
}

func (h *WorkspaceHandler) UpdateWorkspace(c *gin.Context) {
	id := c.Param("id")
	workspaceID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	var updateData models.Workspace
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	workspace, err := h.workspaceRepo.GetByID(workspaceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
		return
	}

	// Check ownership
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user := currentUser.(*models.User)

	if workspace.OwnerID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this workspace"})
		return
	}

	workspace.Name = updateData.Name
	// Add other updateable fields as needed

	if err := h.workspaceRepo.Update(workspace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workspace"})
		return
	}

	c.JSON(http.StatusOK, workspace)
}

func (h *WorkspaceHandler) DeleteWorkspace(c *gin.Context) {
	id := c.Param("id")
	workspaceID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	workspace, err := h.workspaceRepo.GetByID(workspaceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
		return
	}

	// Check ownership
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user := currentUser.(*models.User)

	if workspace.OwnerID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this workspace"})
		return
	}

	if err := h.workspaceRepo.Delete(workspaceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workspace"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workspace deleted successfully"})
}

func (h *WorkspaceHandler) SwitchWorkspace(c *gin.Context) {
	id := c.Param("id")
	workspaceID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	// Check if user has access to the workspace
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user := currentUser.(*models.User)

	// Verify workspace exists and user is the owner
	workspace, err := h.workspaceRepo.GetByID(workspaceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
		return
	}

	if workspace.OwnerID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this workspace"})
		return
	}

	// Update user's current workspace
	user.CurrentWorkspaceID = workspaceID
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to switch workspace"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workspace switched successfully"})
}*/
//package handlers
/*
import (
	"net/http"
	"context"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hashmi846003/P.M.S/internal/models"
	"github.com/hashmi846003/P.M.S/internal/repository"
)

type WorkspaceHandler struct {
	workspaceRepo repository.WorkspaceRepository
	userRepo      repository.UserRepository
}

func NewWorkspaceHandler(workspaceRepo repository.WorkspaceRepository, userRepo repository.UserRepository) *WorkspaceHandler {
	return &WorkspaceHandler{
		workspaceRepo: workspaceRepo,
		userRepo:      userRepo,
	}
}

func (h *WorkspaceHandler) CreateWorkspace(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Get current user from context
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Create workspace
	workspace := &models.Workspace{
		ID:        uuid.New(),
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	// Create workspace
	if err := h.workspaceRepo.Create(c.Request.Context(), workspace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workspace"})
		return
	}

	// Add user to workspace as owner
	if err := h.workspaceRepo.AddUserToWorkspace(
		c.Request.Context(),
		workspace.ID,
		userID,
		true, // isOwner
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to workspace"})
		return
	}

	c.JSON(http.StatusCreated, workspace)
}

func (h *WorkspaceHandler) ListWorkspaces(c *gin.Context) {
	// Get current user
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	workspaces, err := h.workspaceRepo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch workspaces"})
		return
	}
	
	c.JSON(http.StatusOK, workspaces)
}

func (h *WorkspaceHandler) UpdateWorkspace(c *gin.Context) {
	id := c.Param("id")
	workspaceID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	// Get current user
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user is owner of workspace
	isOwner, err := h.workspaceRepo.IsWorkspaceOwner(c.Request.Context(), workspaceID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ownership"})
		return
	}
	
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this workspace"})
		return
	}

	var updateData struct {
		Name string `json:"name" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	workspace, err := h.workspaceRepo.GetByID(c.Request.Context(), workspaceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
		return
	}

	workspace.Name = updateData.Name
	workspace.UpdatedAt = time.Now()

	if err := h.workspaceRepo.Update(c.Request.Context(), workspace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workspace"})
		return
	}

	c.JSON(http.StatusOK, workspace)
}

func (h *WorkspaceHandler) DeleteWorkspace(c *gin.Context) {
	id := c.Param("id")
	workspaceID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	// Get current user
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user is owner of workspace
	isOwner, err := h.workspaceRepo.IsWorkspaceOwner(c.Request.Context(), workspaceID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ownership"})
		return
	}
	
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this workspace"})
		return
	}

	if err := h.workspaceRepo.Delete(c.Request.Context(), workspaceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workspace"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workspace deleted successfully"})
}

func (h *WorkspaceHandler) SwitchWorkspace(c *gin.Context) {
	id := c.Param("id")
	workspaceID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	// Get current user
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user has access to the workspace
	hasAccess, err := h.workspaceRepo.UserHasAccess(c.Request.Context(), workspaceID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify workspace access"})
		return
	}
	
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this workspace"})
		return
	}

	// Update user's current workspace
	if err := h.userRepo.SetCurrentWorkspace(c.Request.Context(), userID, workspaceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to switch workspace"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workspace switched successfully"})
}

// Helper function
func getUserUUID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userId")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return uuid.Parse(userID.(string))
}*/
package handlers

import (
	"net/http"
	//"context"
	//"fmt"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hashmi846003/P.M.S/internal/models"
	"github.com/hashmi846003/P.M.S/internal/repository"
)

type WorkspaceHandler struct {
	workspaceRepo repository.WorkspaceRepository
	userRepo      repository.UserRepository
}

func NewWorkspaceHandler(
	workspaceRepo repository.WorkspaceRepository, 
	userRepo repository.UserRepository,
) *WorkspaceHandler {
	return &WorkspaceHandler{
		workspaceRepo: workspaceRepo,
		userRepo:      userRepo,
	}
}

func (h *WorkspaceHandler) CreateWorkspace(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Get current user
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Create workspace
	workspace := &models.Workspace{
		ID:        uuid.New(),
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	// Create workspace
	if err := h.workspaceRepo.Create(c.Request.Context(), workspace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workspace"})
		return
	}

	// Add user to workspace as owner
	if err := h.workspaceRepo.AddUserToWorkspace(
		c.Request.Context(),
		workspace.ID,
		userID,
		true, // isOwner
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to workspace"})
		return
	}

	c.JSON(http.StatusCreated, workspace)
}

func (h *WorkspaceHandler) ListWorkspaces(c *gin.Context) {
	// Get current user
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	workspaces, err := h.workspaceRepo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch workspaces"})
		return
	}
	
	c.JSON(http.StatusOK, workspaces)
}

func (h *WorkspaceHandler) UpdateWorkspace(c *gin.Context) {
	id := c.Param("id")
	workspaceID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	// Get current user
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user is owner of workspace
	isOwner, err := h.workspaceRepo.IsWorkspaceOwner(c.Request.Context(), workspaceID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ownership"})
		return
	}
	
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this workspace"})
		return
	}

	var updateData struct {
		Name string `json:"name" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	workspace, err := h.workspaceRepo.GetByID(c.Request.Context(), workspaceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
		return
	}
	if workspace == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
		return
	}

	workspace.Name = updateData.Name
	workspace.UpdatedAt = time.Now()

	if err := h.workspaceRepo.Update(c.Request.Context(), workspace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workspace"})
		return
	}

	c.JSON(http.StatusOK, workspace)
}

func (h *WorkspaceHandler) DeleteWorkspace(c *gin.Context) {
	id := c.Param("id")
	workspaceID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	// Get current user
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user is owner of workspace
	isOwner, err := h.workspaceRepo.IsWorkspaceOwner(c.Request.Context(), workspaceID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ownership"})
		return
	}
	
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this workspace"})
		return
	}

	if err := h.workspaceRepo.Delete(c.Request.Context(), workspaceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workspace"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workspace deleted successfully"})
}

func (h *WorkspaceHandler) SwitchWorkspace(c *gin.Context) {
	id := c.Param("id")
	workspaceID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	// Get current user
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user has access to the workspace
	hasAccess, err := h.workspaceRepo.UserHasAccess(c.Request.Context(), workspaceID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify workspace access"})
		return
	}
	
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this workspace"})
		return
	}

	// Update user's current workspace
	if err := h.userRepo.SetCurrentWorkspace(c.Request.Context(), userID, workspaceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to switch workspace"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workspace switched successfully"})
}
/*
// Helper function
func getUserUUID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userId")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return uuid.Parse(userID.(string))
}*/