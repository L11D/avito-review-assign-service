package handlers

import (
	"context"

	"github.com/L11D/avito-review-assign-service/internal/errors"
	"github.com/L11D/avito-review-assign-service/internal/http/dto"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	SetIsActive(ctx context.Context, userSetIsActiveDTO dto.UserSetIsActiveDTO) (dto.UserDTO, error)
	GetReviews(ctx context.Context, userId string) (dto.UserPRsDTO, error)
}

type UserHandler struct{
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) RegisterRoutes(e *gin.Engine) {
	g := e.Group("/users")
	g.POST("/setIsActive", h.setIsActive)
	g.GET("/getReview", h.getReviews)
}

func (h *UserHandler) setIsActive (c *gin.Context){
	var dto dto.UserSetIsActiveDTO
	if err := c.ShouldBindJSON(&dto); err != nil { 
		c.Error(errors.NewValidationFailedError(err.Error()))
		return
	}
	updatedUser, err := h.service.SetIsActive(c.Request.Context(), dto)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, updatedUser)
}

func (h *UserHandler) getReviews (c *gin.Context){
	userId := c.Query("user_id")
	if userId == "" {
		c.Error(errors.NewValidationFailedError("user_id is required"))
		return
	}
	userPRs, err := h.service.GetReviews(c.Request.Context(), userId)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, userPRs)
}