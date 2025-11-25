package app

import (
	"context"
	"errors"
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

	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("Failed to close database connection", slog.String("error", err.Error()))
		}
	}()

	if err := migrations.RunMigrations(db.DB); err != nil {
		slog.Error("Failed to run migrations", slog.String("error", err.Error()))

		return
	}

	r := initDependencies(db)

	server := &http.Server{
		Addr:              ":" + config.HTTPPort,
		Handler:           r,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
	}

	go func() {
		slog.Info("Starting server on :" + config.HTTPPort)

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
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

func initDependencies(db *sqlx.DB) *gin.Engine {
	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))

	userRepo := repo.NewUserRepo(db, trmsqlx.DefaultCtxGetter)
	teamRepo := repo.NewTeamRepo(db, trmsqlx.DefaultCtxGetter)
	pullRequestRepo := repo.NewPullRequestRepo(db, trmsqlx.DefaultCtxGetter)
	pullRequestReviewerRepo := repo.NewPullRequestReviewerRepo(db, trmsqlx.DefaultCtxGetter)

	userService := services.NewUserService(userRepo, teamRepo, pullRequestRepo, trManager)
	teamService := services.NewTeamService(teamRepo, userService, trManager)
	pullService := services.NewPullRequestService(
		pullRequestRepo,
		pullRequestReviewerRepo,
		userRepo,
		userService,
		trManager,
	)
	statisticService := services.NewStatisticService(userRepo, trManager)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorMiddleware())
	r.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"status": "healthy"}) })

	handlers.NewUserHandler(userService).RegisterRoutes(r)
	handlers.NewTeamHandler(teamService).RegisterRoutes(r)
	handlers.NewPullRequestHandler(pullService).RegisterRoutes(r)
	handlers.NewStatisticHandler(statisticService).RegisterRoutes(r)

	return r
}
