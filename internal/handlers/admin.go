
package handlers

import (
	"net/http"
	//"context"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hashmi846003/P.M.S/internal/repository"
)

type AdminHandler struct {
	userRepo repository.UserRepository
}

func NewAdminHandler(userRepo repository.UserRepository) *AdminHandler {
	return &AdminHandler{userRepo: userRepo}
}

func (h *AdminHandler) ListPendingUsers(c *gin.Context) {
	// Pass request context to repository
	users, err := h.userRepo.GetPendingUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *AdminHandler) ApproveUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Pass request context to repository
	if err := h.userRepo.ApproveUser(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve user"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "User approved successfully"})
}

func (h *AdminHandler) RejectUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Pass request context to repository
	if err := h.userRepo.RejectUser(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject user"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "User rejected successfully"})
}