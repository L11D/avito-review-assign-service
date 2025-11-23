package app

import (
	"github.com/gin-gonic/gin"
	"github.com/L11D/avito-review-assign-service/internal/http/handlers"
)

func Run() {
	r := gin.Default()
	teamHandler := handlers.NewTeamHandler()
	teamHandler.RegisterRoutes(r)

	r.Run(":8080")
}