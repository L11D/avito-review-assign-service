package handlers

import (
	"context"
	"net/http"

	"github.com/L11D/avito-review-assign-service/internal/http/dto"
	"github.com/gin-gonic/gin"
)

type TeamService interface {
	Create(ctx context.Context, team dto.TeamDTO) (dto.TeamDTO, error)
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
	if err := c.ShouldBindJSON(&dto); err != nil { // временное решение
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTeam, err := h.service.Create(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdTeam)
}

func (h *TeamHandler) Get(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"status": "team retrieved", "id": id})
}