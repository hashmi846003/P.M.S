package handlers

import (
	//"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"internal/models"
)

type PageHandler struct {
	DB *gorm.DB
}

func NewPageHandler(db *gorm.DB) *PageHandler {
	return &PageHandler{DB: db}
}

// Authentication middleware for page handlers
func (h *PageHandler) authenticateUser(c *gin.Context) (*models.User, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return nil, false
	}

	var user models.User
	if err := h.DB.First(&user, userID.(uint)).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user credentials"})
		return nil, false
	}

	if !user.AdminApproval {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account pending admin approval"})
		return nil, false
	}

	return &user, true
}

// CreatePage creates a new document page
func (h *PageHandler) CreatePage(c *gin.Context) {
	user, ok := h.authenticateUser(c)
	if !ok {
		return
	}

	var input struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page := models.Page{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: user.ID,
	}

	if err := h.DB.Create(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create page"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      page.ID,
		"title":   page.Title,
		"message": "Page created successfully",
	})
}

// GetPage retrieves a specific page with authorization check
func (h *PageHandler) GetPage(c *gin.Context) {
	user, ok := h.authenticateUser(c)
	if !ok {
		return
	}

	pageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	var page models.Page
	if err := h.DB.Preload("Discussions").
		Where("id = ? AND author_id = ?", pageID, user.ID).
		First(&page).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found or access denied"})
		return
	}

	c.JSON(http.StatusOK, page)
}

// UpdatePage modifies an existing page with version history
func (h *PageHandler) UpdatePage(c *gin.Context) {
	user, ok := h.authenticateUser(c)
	if !ok {
		return
	}

	pageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingPage models.Page
	if err := h.DB.Where("id = ? AND author_id = ?", pageID, user.ID).
		First(&existingPage).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found or access denied"})
		return
	}

	// Save to history
	h.DB.Create(&models.PageHistory{
		PageID:    existingPage.ID,
		Content:   existingPage.Content,
		Version:   getNextVersion(existingPage.ID),
		UpdatedBy: user.ID,
	})

	// Update page
	updates := models.Page{
		Title:   input.Title,
		Content: input.Content,
	}

	if err := h.DB.Model(&existingPage).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      existingPage.ID,
		"title":   existingPage.Title,
		"message": "Page updated successfully",
	})
}

// DeletePage handles page deletion with authorization
func (h *PageHandler) DeletePage(c *gin.Context) {
	user, ok := h.authenticateUser(c)
	if !ok {
		return
	}

	pageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	result := h.DB.Where("id = ? AND author_id = ?", pageID, user.ID).Delete(&models.Page{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete page"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found or access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}

// Helper function to get next version number
func getNextVersion(pageID uint) int {
	var count int64
	h.DB.Model(&models.PageHistory{}).Where("page_id = ?", pageID).Count(&count)
	return int(count) + 1
}