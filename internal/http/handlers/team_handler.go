package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/L11D/avito-review-assign-service/internal/http/dto"
)

type TeamHandler struct{}

func NewTeamHandler() *TeamHandler {
	return &TeamHandler{}
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
	c.JSON(http.StatusOK, gin.H{"status": "team added"})
}

func (h *TeamHandler) Get(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"status": "team retrieved", "id": id})
}