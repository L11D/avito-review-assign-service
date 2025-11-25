package handlers

import (
	"context"

	"github.com/L11D/avito-review-assign-service/pkg/api/dto"
	"github.com/gin-gonic/gin"
)

type StatisticService interface {
	GetUsersStatistic(ctx context.Context) (dto.AllUsersStatisticDTO, error)
}

type StatisticHandler struct {
	service StatisticService
}

func NewStatisticHandler(service StatisticService) *StatisticHandler {
	return &StatisticHandler{
		service: service,
	}
}

func (h *StatisticHandler) RegisterRoutes(e *gin.Engine) {
	g := e.Group("/statistic")
	g.GET("/users", h.GetUsersStatistic)
}

func (h *StatisticHandler) GetUsersStatistic(c *gin.Context) {
	statistic, err := h.service.GetUsersStatistic(c.Request.Context())
	if err != nil {
		c.Error(err)

		return
	}

	c.JSON(200, statistic)
}
