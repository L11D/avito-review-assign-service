package app

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	
	"github.com/gin-gonic/gin"
	"github.com/L11D/avito-review-assign-service/internal/http/handlers"
	"github.com/L11D/avito-review-assign-service/internal/services"
	"github.com/L11D/avito-review-assign-service/internal/repo"
	"fmt"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

func Run() {

	db, err := sqlx.Connect("postgres", "postgres://user:pass@localhost:5432/review-assign-service-db?sslmode=disable")
    if err != nil {
        fmt.Println("Failed to connect to the database:", err)
        return
    }

	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))
	
	r := gin.Default()

	userRepo := repo.NewUserRepo(db, trmsqlx.DefaultCtxGetter)
	userService := services.NewUserService(userRepo)

	teamRepo := repo.NewTeamRepo(db, trmsqlx.DefaultCtxGetter)
	teamService := services.NewTeamService(teamRepo, userService, trManager)
	teamHandler := handlers.NewTeamHandler(teamService)


	teamHandler.RegisterRoutes(r)

	r.Run(":8080")
}