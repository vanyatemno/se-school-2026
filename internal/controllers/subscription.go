package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"se-school/internal/models/dto"
	"se-school/internal/services/subscription"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SubscriptionController struct {
	subscriptionService subscription.SubscriptionsService
}

// NewSubscriptionController creates a new SubscriptionController backed by the given service.
func NewSubscriptionController(
	subscriptionService subscription.SubscriptionsService,
) *SubscriptionController {
	return &SubscriptionController{
		subscriptionService: subscriptionService,
	}
}

// Subscribe handles POST /api/subscribe.
// Accepts form-data or JSON with "email" and "repo" fields.
func (sc *SubscriptionController) Subscribe(c *gin.Context) {
	var req dto.CreateSubscriptionRequest

	// Read the raw body so we can attempt multiple parse strategies.
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zap.L().Warn("failed to read request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Try JSON first.
	if json.Unmarshal(body, &req) != nil || req.Email == "" || req.Repo == "" {
		// Reset and try form-urlencoded parsing.
		req = dto.CreateSubscriptionRequest{}
		values, parseErr := url.ParseQuery(string(body))
		if parseErr != nil {
			zap.L().Warn("invalid subscribe request: cannot parse body")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		req.Email = values.Get("email")
		req.Repo = values.Get("repo")
	}

	if req.Email == "" || req.Repo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = sc.subscriptionService.Create(c.Request.Context(), &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription successful. Confirmation email sent."})
}

// Confirm handles GET /api/confirm/:token.
func (sc *SubscriptionController) Confirm(c *gin.Context) {
	var req dto.ConfirmSubscriptionRequest
	if err := c.ShouldBindUri(&req); err != nil {
		zap.L().Warn("invalid confirm request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	err := sc.subscriptionService.Confirm(&req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription confirmed successfully"})
}

// Unsubscribe handles GET /api/unsubscribe/:token.
func (sc *SubscriptionController) Unsubscribe(c *gin.Context) {
	var req dto.UnsubscribeRequest
	if err := c.ShouldBindUri(&req); err != nil {
		zap.L().Warn("invalid unsubscribe request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	err := sc.subscriptionService.Unsubscribe(&req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}

// GetSubscriptions handles GET /api/subscriptions?email=...
func (sc *SubscriptionController) GetSubscriptions(c *gin.Context) {
	var req dto.GetSubscriptionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		zap.L().Warn("invalid get subscriptions request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	subscriptions, err := sc.subscriptionService.ListByEmail(&req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response := make([]dto.SubscriptionResponse, 0, len(subscriptions))
	for _, sub := range subscriptions {
		repo := ""
		if sub.Repository != nil {
			repo = sub.Repository.Owner + "/" + sub.Repository.Name
		}
		response = append(response, dto.SubscriptionResponse{
			Email:       sub.Email,
			Repo:        repo,
			Confirmed:   sub.IsConfirmed,
			LastSeenTag: sub.LastSeenTag,
		})
	}

	c.JSON(http.StatusOK, response)
}
