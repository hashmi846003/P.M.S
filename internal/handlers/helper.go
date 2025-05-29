package handlers

import (
	"fmt"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Common helper function for all handlers
func getUserUUID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userId")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return uuid.Parse(userID.(string))
}