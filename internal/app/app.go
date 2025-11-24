package app

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/L11D/avito-review-assign-service/internal/http/handlers"
	"github.com/L11D/avito-review-assign-service/internal/http/middleware"
	"github.com/L11D/avito-review-assign-service/internal/repo"
	"github.com/L11D/avito-review-assign-service/internal/services"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gin-gonic/gin"
)

func Run() {

	db, err := sqlx.Connect("postgres", "postgres://user:pass@localhost:5432/review-assign-service-db?sslmode=disable")
    if err != nil {
		slog.Error("Failed to connect to the database", slog.String("error", err.Error()))
        return
    }

	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))
	
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandler())
	
	userRepo := repo.NewUserRepo(db, trmsqlx.DefaultCtxGetter)
	teamRepo := repo.NewTeamRepo(db, trmsqlx.DefaultCtxGetter)

	userService := services.NewUserService(userRepo, teamRepo)
	teamService := services.NewTeamService(teamRepo, userService, trManager)

	teamHandler := handlers.NewTeamHandler(teamService)
	userHandler := handlers.NewUserHandler(userService)

	userHandler.RegisterRoutes(r)
	teamHandler.RegisterRoutes(r)

	r.Run(":8080")
}