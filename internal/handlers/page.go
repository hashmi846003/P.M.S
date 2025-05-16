package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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

// GetAllPages returns all pages for the authenticated user
func (h *PageHandler) GetAllPages(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	
	var pages []models.Page
	h.DB.Where("author_id = ? AND is_trash = false", userID).Find(&pages)
	
	c.JSON(http.StatusOK, pages)
}

// CreatePage handles new page creation
func (h *PageHandler) CreatePage(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page := models.Page{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: userID,
	}

	if err := h.DB.Create(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create page"})
		return
	}

	c.JSON(http.StatusCreated, page)
}

// GetPage retrieves a single page with discussions
func (h *PageHandler) GetPage(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	pageID := parseUint(c.Param("id"))

	var page models.Page
	if err := h.DB.Preload("Discussions").
		Where("id = ? AND author_id = ?", pageID, userID).
		First(&page).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	c.JSON(http.StatusOK, page)
}

// UpdatePage modifies an existing page and saves history
func (h *PageHandler) UpdatePage(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	pageID := parseUint(c.Param("id"))

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var page models.Page
	if err := h.DB.Where("id = ? AND author_id = ?", pageID, userID).First(&page).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	// Save to history
	h.DB.Create(&models.PageHistory{
		PageID:    page.ID,
		Content:   page.Content,
		Version:   getNextVersion(page.ID),
		UpdatedBy: userID,
	})

	page.Title = input.Title
	page.Content = input.Content
	h.DB.Save(&page)

	c.JSON(http.StatusOK, page)
}

// DeletePage soft deletes a page
func (h *PageHandler) DeletePage(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	pageID := parseUint(c.Param("id"))
