package codes

import "se-school/internal/models"

type CodesService interface {
	Get(code string) (*models.Code, error)
	Create(subscriptionID uint, codeType models.CodeType) (*models.Code, error)
	Delete(id uint) error
}
