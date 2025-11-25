package handlers

import (
	"context"

	"github.com/L11D/avito-review-assign-service/internal/errors"
	"github.com/L11D/avito-review-assign-service/internal/http/dto"
	"github.com/gin-gonic/gin"
)

type PullRequestService interface {
	Create(ctx context.Context, pr dto.PullRequestCreateDTO) (dto.PullRequestDTO, error)
	Merge(ctx context.Context, prId string) (dto.PullRequestDTO, error)
	Reassign(ctx context.Context, reassignDTO dto.PullRequestReassignDTO) (dto.PullRequestDTO, error)
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
	g.POST("/merge", h.Merge)
	g.POST("/reassign", h.Reassign)
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

func (h *PullRequestHandler) Merge(c *gin.Context) {
	var dto dto.PullRequestMergeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(errors.NewValidationFailedError(err.Error()))
		return
	}

	mergedPR,  err := h.service.Merge(c.Request.Context(), dto.Id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"pr": mergedPR})
}

func (h *PullRequestHandler) Reassign(c *gin.Context) {
	var dto dto.PullRequestReassignDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(errors.NewValidationFailedError(err.Error()))
		return
	}

	reassignedPR,  err := h.service.Reassign(c.Request.Context(), dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"pr": reassignedPR})
}