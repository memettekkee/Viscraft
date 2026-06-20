package handler

import (
	"net/http"

	"viscraft-backend/repository"

	"github.com/gin-gonic/gin"
)

type PromptOptionController struct {
	repo *repository.PromptOptionRepository
}

func NewPromptOptionController(repo *repository.PromptOptionRepository) *PromptOptionController {
	return &PromptOptionController{repo: repo}
}

func (c *PromptOptionController) ListByCategory(ctx *gin.Context) {
	options, err := c.repo.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch options"})
		return
	}

	if options == nil {
		options = []repository.PromptOption{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    options,
	})
}
