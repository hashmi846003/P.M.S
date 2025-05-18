package handlers

import (
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	"your-app/internal/models"
	"your-app/internal/repository"
)

type PageHandler struct {
	repo repository.PageRepository
}

func NewPageHandler(repo repository.PageRepository) *PageHandler {
	return &PageHandler{repo: repo}
}

func (h *PageHandler) CreatePage(c *gin.Context) {
	userID, _ := c.Get("userID")
	var page models.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	page.UserID = userID.(uint)
	createdPage, err := h.repo.CreatePage(c.Request.Context(), &page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create page"})
		return
	}

	c.JSON(http.StatusCreated, createdPage)
}

func (h *PageHandler) GetAllPages(c *gin.Context) {
	userID, _ := c.Get("userID")
	pages, err := h.repo.GetAllPages(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pages"})
		return
	}
	c.JSON(http.StatusOK, pages)
}

func (h *PageHandler) GetPageByID(c *gin.Context) {
	userID, _ := c.Get("userID")
	pageID, _ := strconv.Atoi(c.Param("id"))
	
	page, err := h.repo.GetPageByID(c.Request.Context(), uint(pageID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	if page.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) UpdatePage(c *gin.Context) {
	userID, _ := c.Get("userID")
	pageID, _ := strconv.Atoi(c.Param("id"))

	var page models.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	existingPage, err := h.repo.GetPageByID(c.Request.Context(), uint(pageID))
	if err != nil || existingPage.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	page.ID = uint(pageID)
	page.UserID = userID.(uint)
	updatedPage, err := h.repo.UpdatePage(c.Request.Context(), &page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update page"})
		return
	}

	c.JSON(http.StatusOK, updatedPage)
}

func (h *PageHandler) DeletePage(c *gin.Context) {
	userID, _ := c.Get("userID")
	pageID, _ := strconv.Atoi(c.Param("id"))

	page, err := h.repo.GetPageByID(c.Request.Context(), uint(pageID))
	if err != nil || page.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.repo.DeletePage(c.Request.Context(), uint(pageID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}