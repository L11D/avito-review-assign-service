package handlers

import (
	"context"
	"net/http"

	"github.com/L11D/avito-review-assign-service/internal/errors"
	"github.com/L11D/avito-review-assign-service/internal/http/dto"
	"github.com/gin-gonic/gin"
)

type TeamService interface {
	Create(ctx context.Context, team dto.TeamDTO) (dto.TeamDTO, error)
	GetByName(ctx context.Context, name string) (dto.TeamDTO, error)
}

type TeamHandler struct{
	service TeamService
}

func NewTeamHandler(service TeamService) *TeamHandler {
	return &TeamHandler{
		service: service,
	}
}

func (h *TeamHandler) RegisterRoutes(e *gin.Engine) {
	g := e.Group("/team")
	g.POST("/add", h.Add)
	g.GET("/get", h.Get)
}

func (h *TeamHandler) Add(c *gin.Context) {
	var dto dto.TeamDTO
	if err := c.ShouldBindJSON(&dto); err != nil { 
		c.Error(errors.NewValidationFailedError(err.Error()))
		return
	}

	createdTeam, err := h.service.Create(c.Request.Context(), dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, createdTeam)
}

func (h *TeamHandler) Get(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.Error(errors.NewQueryParamMissingError("name"))
		return
	}
	team, err := h.service.GetByName(c.Request.Context(), name)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, team)
}