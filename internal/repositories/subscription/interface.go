package subscription

import "se-school/internal/models"

type SubscriptionsRepository interface {
	GetByID(uint) (*models.Subscription, error)
	GetByRepositoryID(uint) ([]*models.Subscription, error)
	GetByCode(codeID uint, codeType models.CodeType) (*models.Subscription, error)
	GetByEmail(string) ([]*models.Subscription, error)
	Create(subscription *models.Subscription) error
	UpdateLastSeenTag(id uint, tag string) error
	Save(*models.Subscription) error
	Delete(uint) error
}
