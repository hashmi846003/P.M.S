// internal/handlers/emoji.go
package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/enescakir/emoji"
)

type EmojiHandler struct{}

func NewEmojiHandler() *EmojiHandler {
	return &EmojiHandler{}
}

// ListEmojis returns all available emojis
func (h *EmojiHandler) ListEmojis(c *gin.Context) {
	emojis := make(map[string][]string)
	
	// Organize emojis by category
	emojis["Smileys"] = []string{
		emoji.GrinningFace.String(),
		emoji.BeamingFaceWithSmilingEyes.String(),
		emoji.FaceWithTearsOfJoy.String(),
		// Add more from emoji package...
	}
	
	emojis["Nature"] = []string{
		emoji.Dog.String(),
		emoji.Cat.String(),
		//emoji.Tree.String(),
		// Add more...
	}
	
	emojis["Foods"] = []string{
		emoji.Pizza.String(),
		emoji.Hamburger.String(),
		emoji.Taco.String(),
		// Add more...
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": emojis,
	})
}

// GetCategories returns emoji categories
func (h *EmojiHandler) GetCategories(c *gin.Context) {
	categories := []string{
		"Smileys & Emotion",
		"People & Body",
		"Animals & Nature",
		"Food & Drink",
		"Travel & Places",
		"Activities",
		"Objects",
		"Symbols",
		"Flags",
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}