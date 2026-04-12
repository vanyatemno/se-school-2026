package main

import (
	"log"
	"se-school/internal/notifications"
	"se-school/internal/notifications/mailer"
	"se-school/internal/notifications/templates"

	"se-school/internal/config"
	"se-school/internal/controllers"
	"se-school/internal/db"
	"se-school/internal/integrations/github"
	codeRepo "se-school/internal/repositories/code"
	repoRepo "se-school/internal/repositories/repository"
	subRepo "se-school/internal/repositories/subscription"
	subscriptionSvc "se-school/internal/services/subscription"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	zap.ReplaceGlobals(logger)

	cfg, err := config.Read()
	if err != nil {
		zap.L().Fatal("failed to read config", zap.Error(err))
	}

	database, err := db.Connect(&cfg.Database)
	if err != nil {
		zap.L().Fatal("failed to connect to database", zap.Error(err))
	}

	// Repositories
	subscriptionRepository := subRepo.New(database)
	repositoryRepository := repoRepo.New(database)
	codeRepository := codeRepo.New(database)

	// Integrations
	githubIntegration := github.New(&cfg.Github)

	// Services
	subscriptionService := subscriptionSvc.New(
		subscriptionRepository,
		repositoryRepository,
		codeRepository,
		githubIntegration,
		notifications.New(
			mailer.NewMailerService(&cfg.Mailer),
			templates.New(),
		),
	)

	// Controllers
	subscriptionController := controllers.NewSubscriptionController(subscriptionService)

	// Router
	r := gin.Default()
	controllers.RegisterRoutes(r, subscriptionController)

	port := cfg.Application.Port
	if port == "" {
		port = "8080"
	}

	zap.L().Info("starting server", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		zap.L().Fatal("failed to start server", zap.Error(err))
	}
}
