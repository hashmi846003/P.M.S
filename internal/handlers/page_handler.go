/*package handlers

import (
	"net/http"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hashmi846003/P.M.S/internal/models"
	"github.com/hashmi846003/P.M.S/internal/repository"
)

type PageHandler struct {
	pageRepo       repository.PageRepository
	discussionRepo repository.DiscussionRepository
}
// Add these methods to your existing PageHandler struct

func (h *PageHandler) FormatContent(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var req struct {
		Format string `json:"format" binding:"required"`
		Start  int    `json:"start" binding:"required"`
		End    int    `json:"end" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	page, err := h.pageRepo.ApplyFormatting(c.Request.Context(), pageID, userID, req.Format, req.Start, req.End)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply formatting"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) AlignText(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var req struct {
		Alignment string `json:"alignment" binding:"required,oneof=left right center justify"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	page, err := h.pageRepo.SetAlignment(c.Request.Context(), pageID, userID, req.Alignment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set alignment"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) AddEmoji(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var req struct {
		Emoji string `json:"emoji" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	page, err := h.pageRepo.AddEmoji(c.Request.Context(), pageID, userID, req.Emoji)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add emoji"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) GenerateShareLink(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	// Generate a random token (you might want to use a proper token generation library)
	token := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 1 week expiration

	if err := h.pageRepo.CreateShareLink(c.Request.Context(), pageID, token, expiresAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate share link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"share_url":   "/shared/" + token,
		"expires_at": expiresAt.Format(time.RFC3339),
	})
}

func (h *PageHandler) ListTrash(c *gin.Context) {
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	pages, err := h.pageRepo.FindByUserID(c.Request.Context(), userID, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list trash"})
		return
	}

	// Filter only deleted items
	var trash []models.Page
	for _, page := range pages {
		if page.IsDeleted {
			trash = append(trash, page)
		}
	}

	c.JSON(http.StatusOK, trash)
}

func (h *PageHandler) MoveToTrash(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	if err := h.pageRepo.SoftDelete(c.Request.Context(), pageID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to move to trash"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page moved to trash successfully"})
}

func NewPageHandler(pr repository.PageRepository, dr repository.DiscussionRepository) *PageHandler {
	return &PageHandler{
		pageRepo:       pr,
		discussionRepo: dr,
	}
}

func (h *PageHandler) ListPages(c *gin.Context) {
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	includeDeleted := c.Query("deleted") == "true"
	pages, err := h.pageRepo.FindByUserID(c.Request.Context(), userID, includeDeleted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pages"})
		return
	}

	c.JSON(http.StatusOK, pages)
}

func (h *PageHandler) CreatePage(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	newPage := &models.Page{
		ID:        uuid.New(),
		Title:     req.Title,
		Content:   req.Content,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.pageRepo.Create(c.Request.Context(), newPage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create page"})
		return
	}

	c.JSON(http.StatusCreated, newPage)
}

func (h *PageHandler) GetPage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	page, err := h.pageRepo.FindByID(c.Request.Context(), pageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) UpdatePage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var update struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	page, err := h.pageRepo.FindByID(c.Request.Context(), pageID)
	if err != nil || page.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	page.Title = update.Title
	page.Content = update.Content
	page.UpdatedAt = time.Now()

	if err := h.pageRepo.Update(c.Request.Context(), page); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update page"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) DeletePage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	if err := h.pageRepo.SoftDelete(c.Request.Context(), pageID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}

func (h *PageHandler) RestorePage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	if err := h.pageRepo.Restore(c.Request.Context(), pageID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page restored successfully"})
}

func (h *PageHandler) ToggleFavorite(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	isFavorite, err := h.pageRepo.ToggleFavorite(c.Request.Context(), pageID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_favorite": isFavorite})
}

func (h *PageHandler) DuplicatePage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	newPage, err := h.pageRepo.Duplicate(c.Request.Context(), pageID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to duplicate page"})
		return
	}

	c.JSON(http.StatusCreated, newPage)
}

func (h *PageHandler) GetVersions(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	versions, err := h.pageRepo.GetVersions(c.Request.Context(), pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get versions"})
		return
	}

	c.JSON(http.StatusOK, versions)
}

func (h *PageHandler) CreateDiscussion(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	discussion := &models.Discussion{
		ID:        uuid.New(),
		PageID:    pageID,
		Content:   req.Content,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := h.discussionRepo.Create(c.Request.Context(), discussion); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create discussion"})
		return
	}

	c.JSON(http.StatusCreated, discussion)
}

func (h *PageHandler) GetDiscussions(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	discussions, err := h.discussionRepo.GetByPageID(c.Request.Context(), pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get discussions"})
		return
	}

	c.JSON(http.StatusOK, discussions)
}

// Helper functions
func getUserUUID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userId")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return uuid.Parse(userID.(string))
}
// File: internal/handlers/page.go (add these methods to PageHandler)
// Add to existing PageHandler struct

func (h *PageHandler) PermanentDelete(c *gin.Context) {
	id := c.Param("id")
	pageID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	// Check if page exists and is in trash
	page, err := h.pageRepo.GetByID(pageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	if !page.IsDeleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page is not in trash. Move to trash first."})
		return
	}

	// Permanent delete
	if err := h.pageRepo.PermanentDelete(pageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to permanently delete page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page permanently deleted"})
}

// Updated method with permission check
func (h *PageHandler) GenerateShareLink(c *gin.Context) {
	id := c.Param("id")
	pageID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	// Get permission from query parameter
	permission := c.Query("permission")
	if permission == "" {
		permission = "viewer" // Default permission
	}
	
	if permission != "viewer" && permission != "editor" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission type. Must be 'viewer' or 'editor'"})
		return
	}

	// Get current user
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user := currentUser.(*models.User)

	// Create share link with permission
	shareLink, err := h.pageRepo.CreateShareLink(pageID, user.ID, permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate share link"})
		return
	}

	c.JSON(http.StatusCreated, shareLink)
}

// Updated method with pagination and sorting
func (h *PageHandler) ListPages(c *gin.Context) {
	// Get current user
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user := currentUser.(*models.User)

	// Get pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	sortBy := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	// Convert to int
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10
	}

	// Get section ID if provided
	sectionIDStr := c.Query("section_id")
	var sectionID *uuid.UUID
	if sectionIDStr != "" {
		id, err := uuid.Parse(sectionIDStr)
		if err == nil {
			sectionID = &id
		}
	}

	// Get pages with pagination
	pages, total, err := h.pageRepo.ListPages(user.CurrentWorkspaceID, sectionID, page, size, sortBy, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pages"})
		return
	}

	response := gin.H{
		"pages": pages,
		"pagination": gin.H{
			"page":  page,
			"size":  size,
			"total": total,
		},
	}
	c.JSON(http.StatusOK, response)
}

//func parseUUIDParam(c *gin.Context, param string) (uuid.UUID, error) {
//	return uuid.Parse(c.Param(param))
//}*/
/*
package handlers

import (
	"net/http"
	"time"
	"fmt"
	"strconv" // Added for pagination

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hashmi846003/P.M.S/internal/models"
	"github.com/hashmi846003/P.M.S/internal/repository"
)

type PageHandler struct {
	pageRepo       repository.PageRepository
	discussionRepo repository.DiscussionRepository
}

// ... (keep all existing methods unchanged until the end) ...

// ========== NEW METHODS ========== //
func (h *PageHandler) FormatContent(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var req struct {
		Format string `json:"format" binding:"required"`
		Start  int    `json:"start" binding:"required"`
		End    int    `json:"end" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	page, err := h.pageRepo.ApplyFormatting(c.Request.Context(), pageID, userID, req.Format, req.Start, req.End)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply formatting"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) AlignText(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var req struct {
		Alignment string `json:"alignment" binding:"required,oneof=left right center justify"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	page, err := h.pageRepo.SetAlignment(c.Request.Context(), pageID, userID, req.Alignment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set alignment"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) AddEmoji(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var req struct {
		Emoji string `json:"emoji" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	page, err := h.pageRepo.AddEmoji(c.Request.Context(), pageID, userID, req.Emoji)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add emoji"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) GenerateShareLink(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	// Generate a random token (you might want to use a proper token generation library)
	token := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 1 week expiration

	if err := h.pageRepo.CreateShareLink(c.Request.Context(), pageID, token, expiresAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate share link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"share_url":   "/shared/" + token,
		"expires_at": expiresAt.Format(time.RFC3339),
	})
}

func (h *PageHandler) ListTrash(c *gin.Context) {
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	pages, err := h.pageRepo.FindByUserID(c.Request.Context(), userID, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list trash"})
		return
	}

	// Filter only deleted items
	var trash []models.Page
	for _, page := range pages {
		if page.IsDeleted {
			trash = append(trash, page)
		}
	}

	c.JSON(http.StatusOK, trash)
}

func (h *PageHandler) MoveToTrash(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	if err := h.pageRepo.SoftDelete(c.Request.Context(), pageID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to move to trash"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page moved to trash successfully"})
}

func NewPageHandler(pr repository.PageRepository, dr repository.DiscussionRepository) *PageHandler {
	return &PageHandler{
		pageRepo:       pr,
		discussionRepo: dr,
	}
}

func (h *PageHandler) ListPages(c *gin.Context) {
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	includeDeleted := c.Query("deleted") == "true"
	pages, err := h.pageRepo.FindByUserID(c.Request.Context(), userID, includeDeleted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pages"})
		return
	}

	c.JSON(http.StatusOK, pages)
}

func (h *PageHandler) CreatePage(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	newPage := &models.Page{
		ID:        uuid.New(),
		Title:     req.Title,
		Content:   req.Content,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.pageRepo.Create(c.Request.Context(), newPage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create page"})
		return
	}

	c.JSON(http.StatusCreated, newPage)
}

func (h *PageHandler) GetPage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	page, err := h.pageRepo.FindByID(c.Request.Context(), pageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) UpdatePage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var update struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	page, err := h.pageRepo.FindByID(c.Request.Context(), pageID)
	if err != nil || page.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	page.Title = update.Title
	page.Content = update.Content
	page.UpdatedAt = time.Now()

	if err := h.pageRepo.Update(c.Request.Context(), page); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update page"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) DeletePage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	if err := h.pageRepo.SoftDelete(c.Request.Context(), pageID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}

func (h *PageHandler) RestorePage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	if err := h.pageRepo.Restore(c.Request.Context(), pageID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page restored successfully"})
}

func (h *PageHandler) ToggleFavorite(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	isFavorite, err := h.pageRepo.ToggleFavorite(c.Request.Context(), pageID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_favorite": isFavorite})
}

func (h *PageHandler) DuplicatePage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	newPage, err := h.pageRepo.Duplicate(c.Request.Context(), pageID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to duplicate page"})
		return
	}

	c.JSON(http.StatusCreated, newPage)
}

func (h *PageHandler) GetVersions(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	versions, err := h.pageRepo.GetVersions(c.Request.Context(), pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get versions"})
		return
	}

	c.JSON(http.StatusOK, versions)
}

func (h *PageHandler) CreateDiscussion(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	discussion := &models.Discussion{
		ID:        uuid.New(),
		PageID:    pageID,
		Content:   req.Content,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := h.discussionRepo.Create(c.Request.Context(), discussion); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create discussion"})
		return
	}

	c.JSON(http.StatusCreated, discussion)
}

func (h *PageHandler) GetDiscussions(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	discussions, err := h.discussionRepo.GetByPageID(c.Request.Context(), pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get discussions"})
		return
	}

	c.JSON(http.StatusOK, discussions)
}

// Helper functions
func getUserUUID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userId")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return uuid.Parse(userID.(string))
}

func (h *PageHandler) PermanentDelete(c *gin.Context) {
	id := c.Param("id")
	pageID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	// Check if page exists and is in trash
	page, err := h.pageRepo.FindByID(c.Request.Context(), pageID) // Changed to FindByID
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	if !page.IsDeleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page is not in trash. Move to trash first."})
		return
	}

	// Permanent delete
	if err := h.pageRepo.PermanentDelete(c.Request.Context(), pageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to permanently delete page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page permanently deleted"})
}

func (h *PageHandler) GenerateShareLink(c *gin.Context) {
	id := c.Param("id")
	pageID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	// Get permission from query parameter
	permission := c.Query("permission")
	if permission == "" {
		permission = "viewer" // Default permission
	}
	
	if permission != "viewer" && permission != "editor" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission type. Must be 'viewer' or 'editor'"})
		return
	}

	// Get current user
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Create share link with permission
	shareLink, err := h.pageRepo.CreateShareLink(
		c.Request.Context(), 
		pageID, 
		userID, 
		permission, 
		time.Now().Add(7*24*time.Hour), // 1 week expiration
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate share link"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"share_url":   "/shared/" + shareLink.Token,
		"expires_at": shareLink.ExpiresAt.Format(time.RFC3339),
		"permission": shareLink.Permission,
	})
}

func (h *PageHandler) ListPages(c *gin.Context) {
	// Get current user
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	sortBy := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	// Convert to int
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10
	}

	// Get section ID if provided
	sectionIDStr := c.Query("section_id")
	var sectionID *uuid.UUID
	if sectionIDStr != "" {
		id, err := uuid.Parse(sectionIDStr)
		if err == nil {
			sectionID = &id
		}
	}

	// Get pages with pagination
	pages, total, err := h.pageRepo.ListPages(
		c.Request.Context(),
		userID,
		sectionID,
		page,
		size,
		sortBy,
		order,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pages"})
		return
	}

	response := gin.H{
		"pages": pages,
		"pagination": gin.H{
			"page":  page,
			"size":  size,
			"total": total,
		},
	}
	c.JSON(http.StatusOK, response)
}

// Helper functions
/*func getUserUUID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userId")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return uuid.Parse(userID.(string))
}*/
//*/
package handlers

import (
	"net/http"
	"time"
	//"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hashmi846003/P.M.S/internal/models"
	"github.com/hashmi846003/P.M.S/internal/repository"
)

type PageHandler struct {
	pageRepo       repository.PageRepository
	discussionRepo repository.DiscussionRepository
}

// ========== FORMATTING METHODS ========== //
func (h *PageHandler) FormatContent(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var req struct {
		Format string `json:"format" binding:"required"`
		Start  int    `json:"start" binding:"required"`
		End    int    `json:"end" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	page, err := h.pageRepo.ApplyFormatting(c.Request.Context(), pageID, userID, req.Format, req.Start, req.End)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply formatting"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) AlignText(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var req struct {
		Alignment string `json:"alignment" binding:"required,oneof=left right center justify"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	page, err := h.pageRepo.SetAlignment(c.Request.Context(), pageID, userID, req.Alignment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set alignment"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) AddEmoji(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var req struct {
		Emoji string `json:"emoji" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	page, err := h.pageRepo.AddEmoji(c.Request.Context(), pageID, userID, req.Emoji)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add emoji"})
		return
	}

	c.JSON(http.StatusOK, page)
}

// ========== SHARING & TRASH METHODS ========== //
func (h *PageHandler) ListTrash(c *gin.Context) {
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	pages, err := h.pageRepo.FindByUserID(c.Request.Context(), userID, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list trash"})
		return
	}

	// Filter only deleted items
	var trash []models.Page
	for _, page := range pages {
		if page.IsDeleted {
			trash = append(trash, page)
		}
	}

	c.JSON(http.StatusOK, trash)
}

func (h *PageHandler) MoveToTrash(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	if err := h.pageRepo.SoftDelete(c.Request.Context(), pageID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to move to trash"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page moved to trash successfully"})
}

func (h *PageHandler) PermanentDelete(c *gin.Context) {
	id := c.Param("id")
	pageID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	// Check if page exists and is in trash
	page, err := h.pageRepo.FindByID(c.Request.Context(), pageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	if !page.IsDeleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page is not in trash. Move to trash first."})
		return
	}

	// Permanent delete
	if err := h.pageRepo.PermanentDelete(c.Request.Context(), pageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to permanently delete page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page permanently deleted"})
}

// ========== PAGE MANAGEMENT METHODS ========== //
func NewPageHandler(pr repository.PageRepository, dr repository.DiscussionRepository) *PageHandler {
	return &PageHandler{
		pageRepo:       pr,
		discussionRepo: dr,
	}
}

func (h *PageHandler) CreatePage(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	newPage := &models.Page{
		ID:        uuid.New(),
		Title:     req.Title,
		Content:   req.Content,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.pageRepo.Create(c.Request.Context(), newPage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create page"})
		return
	}

	c.JSON(http.StatusCreated, newPage)
}

func (h *PageHandler) GetPage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	page, err := h.pageRepo.FindByID(c.Request.Context(), pageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) UpdatePage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var update struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	page, err := h.pageRepo.FindByID(c.Request.Context(), pageID)
	if err != nil || page.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	page.Title = update.Title
	page.Content = update.Content
	page.UpdatedAt = time.Now()

	if err := h.pageRepo.Update(c.Request.Context(), page); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update page"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) DeletePage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	if err := h.pageRepo.SoftDelete(c.Request.Context(), pageID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}

func (h *PageHandler) RestorePage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	if err := h.pageRepo.Restore(c.Request.Context(), pageID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page restored successfully"})
}

func (h *PageHandler) ToggleFavorite(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	isFavorite, err := h.pageRepo.ToggleFavorite(c.Request.Context(), pageID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_favorite": isFavorite})
}

func (h *PageHandler) DuplicatePage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	newPage, err := h.pageRepo.Duplicate(c.Request.Context(), pageID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to duplicate page"})
		return
	}

	c.JSON(http.StatusCreated, newPage)
}

func (h *PageHandler) GetVersions(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	versions, err := h.pageRepo.GetVersions(c.Request.Context(), pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get versions"})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// ========== DISCUSSION METHODS ========== //
func (h *PageHandler) CreateDiscussion(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	discussion := &models.Discussion{
		ID:        uuid.New(),
		PageID:    pageID,
		Content:   req.Content,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := h.discussionRepo.Create(c.Request.Context(), discussion); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create discussion"})
		return
	}

	c.JSON(http.StatusCreated, discussion)
}

func (h *PageHandler) GetDiscussions(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	discussions, err := h.discussionRepo.GetByPageID(c.Request.Context(), pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get discussions"})
		return
	}

	c.JSON(http.StatusOK, discussions)
}

// ========== SHARING METHOD WITH PERMISSIONS ========== //
func (h *PageHandler) GenerateShareLink(c *gin.Context) {
	id := c.Param("id")
	pageID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	// Get permission from query parameter
	permission := c.Query("permission")
	if permission == "" {
		permission = "viewer" // Default permission
	}
	
	if permission != "viewer" && permission != "editor" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission type. Must be 'viewer' or 'editor'"})
		return
	}

	// Get current user
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Create share link with permission
	shareLink, err := h.pageRepo.CreateShareLink(
		c.Request.Context(), 
		pageID, 
		userID, 
		permission, 
		time.Now().Add(7*24*time.Hour), // 1 week expiration
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate share link"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"share_url":   "/shared/" + shareLink.Token,
		"expires_at": shareLink.ExpiresAt.Format(time.RFC3339),
		"permission": shareLink.Permission,
	})
}

// ========== PAGINATED PAGE LISTING ========== //
func (h *PageHandler) ListPages(c *gin.Context) {
	// Get current user
	userID, err := getUserUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	sortBy := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	// Convert to int
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10
	}

	// Get section ID if provided
	sectionIDStr := c.Query("section_id")
	var sectionID *uuid.UUID
	if sectionIDStr != "" {
		id, err := uuid.Parse(sectionIDStr)
		if err == nil {
			sectionID = &id
		}
	}

	// Get pages with pagination
	pages, total, err := h.pageRepo.ListPages(
		c.Request.Context(),
		userID,
		sectionID,
		page,
		size,
		sortBy,
		order,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pages"})
		return
	}

	response := gin.H{
		"pages": pages,
		"pagination": gin.H{
			"page":  page,
			"size":  size,
			"total": total,
		},
	}
	c.JSON(http.StatusOK, response)
}
/*
// Helper functions
func getUserUUID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userId")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return uuid.Parse(userID.(string))
}*/