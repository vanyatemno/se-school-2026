package controllers

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all subscription-related routes on the given Gin engine
// under the /api base path, matching the swagger specification.
// CORS middleware is applied globally to allow cross-origin requests.
func RegisterRoutes(r *gin.Engine, sc *SubscriptionController) {
	r.Use(CORSMiddleware())

	api := r.Group("/api")
	{
		api.POST("/subscribe", sc.Subscribe)
		api.GET("/confirm/:token", sc.Confirm)
		api.GET("/unsubscribe/:token", sc.Unsubscribe)
		api.GET("/subscriptions", sc.GetSubscriptions)
	}
}
