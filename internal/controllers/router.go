package controllers

import (
	"se-school/internal/config"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all subscription-related routes on the given Gin engine
// under the /api base path, matching the swagger specification.
// CORS middleware is applied globally to allow cross-origin requests.
// Protected endpoints require a valid API key via the X-API-Key header.
func RegisterRoutes(r *gin.Engine, sc *SubscriptionController, cfg *config.Application) {
	r.Use(CORSMiddleware())

	api := r.Group("/api", APIKeyMiddleware(cfg.APIKey))
	{
		// Public endpoints — accessed via email confirmation links.
		api.GET("/confirm/:token", sc.Confirm)
		api.GET("/unsubscribe/:token", sc.Unsubscribe)
		api.POST("/subscribe", sc.Subscribe)
		api.GET("/subscriptions", sc.GetSubscriptions)
	}
}
