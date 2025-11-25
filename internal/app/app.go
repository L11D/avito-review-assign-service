package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/L11D/avito-review-assign-service/internal/config"
	"github.com/L11D/avito-review-assign-service/internal/http/handlers"
	"github.com/L11D/avito-review-assign-service/internal/http/middleware"
	"github.com/L11D/avito-review-assign-service/internal/migrations"
	"github.com/L11D/avito-review-assign-service/internal/repo"
	"github.com/L11D/avito-review-assign-service/internal/services"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gin-gonic/gin"
)

func Run() {
	config, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", slog.String("error", err.Error()))
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), 
        os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
    defer stop()

	db, err := sqlx.Connect("postgres", config.GetDBSource())
	if err != nil {
		slog.Error("Failed to connect to the database", slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	if err := migrations.RunMigrations(db.DB); err != nil {
		slog.Error("Failed to run migrations", slog.String("error", err.Error()))
		return
	}

	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))

	userRepo := repo.NewUserRepo(db, trmsqlx.DefaultCtxGetter)
	teamRepo := repo.NewTeamRepo(db, trmsqlx.DefaultCtxGetter)

	userService := services.NewUserService(userRepo, teamRepo, trManager)
	teamService := services.NewTeamService(teamRepo, userService, trManager)

	teamHandler := handlers.NewTeamHandler(teamService)
	userHandler := handlers.NewUserHandler(userService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorMiddleware())

	userHandler.RegisterRoutes(r)
	teamHandler.RegisterRoutes(r)

	server := &http.Server{
        Addr:    ":" + config.HTTPPort,
        Handler: r,
	}

	go func() {
        slog.Info("Starting server on :" + config.HTTPPort)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            slog.Error("Server failed to start", slog.String("error", err.Error()))
            stop() 
        }
    }()

	<-ctx.Done()
	slog.Info("Shutting down application...")
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
    defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
        slog.Error("Server shutdown failed", slog.String("error", err.Error()))
    }

	slog.Info("Application stopped")
}
