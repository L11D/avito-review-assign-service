package app

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	
	"github.com/gin-gonic/gin"
	"github.com/L11D/avito-review-assign-service/internal/http/handlers"
	"github.com/L11D/avito-review-assign-service/internal/services"
	"github.com/L11D/avito-review-assign-service/internal/repo"
	"fmt"
)

func Run() {

	db, err := sqlx.Connect("postgres", "postgres://user:pass@localhost:5432/review-assign-service-db?sslmode=disable")
    if err != nil {
        fmt.Println("Failed to connect to the database:", err)
        return
    }


	r := gin.Default()

	teamRepo := repo.NewTeamRepo(db)
	teamService := services.NewTeamService(teamRepo)
	teamHandler := handlers.NewTeamHandler(teamService)


	teamHandler.RegisterRoutes(r)

	r.Run(":8080")
}