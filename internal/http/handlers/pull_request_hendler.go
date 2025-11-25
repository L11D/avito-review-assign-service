package handlers

import (
	"context"

	"github.com/L11D/avito-review-assign-service/internal/errors"
	"github.com/L11D/avito-review-assign-service/internal/http/dto"
	"github.com/gin-gonic/gin"
)

type PullRequestService interface {
	Create(ctx context.Context, pr dto.PullRequestCreateDTO) (dto.PullRequestDTO, error)
}

type PullRequestHandler struct{
	service PullRequestService
}

func NewPullRequestHandler(service PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{
		service: service,
	}
}

func (h *PullRequestHandler) RegisterRoutes(e *gin.Engine) {
	g := e.Group("/pullRequest")
	g.POST("/create", h.Create)
}

func (h *PullRequestHandler) Create(c *gin.Context) {
	var dto dto.PullRequestCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(errors.NewValidationFailedError(err.Error()))
		return
	}

	createdPR, err := h.service.Create(c.Request.Context(), dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, gin.H{"pr": createdPR})
}