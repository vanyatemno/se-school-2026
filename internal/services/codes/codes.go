package codes

import (
	"se-school/internal/models"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Get(codeString string) (*models.Code, error) {
	var code models.Code
	err := s.db.Where(&models.Code{Code: codeString}).First(&code).Error
	if err != nil {
		return nil, err
	}

	return &code, nil
}

func (s *Service) Create(
	subscriptionID uint,
	codeType models.CodeType,
) (*models.Code, error) {
	code := models.Code{
		SubscriptionID: subscriptionID,
		Code:           "",
		Type:           codeType,
	}
	err := s.setupCode(&code)
	if err != nil {
		return nil, err
	}

	err = s.db.Create(&code).Error
	if err != nil {
		return nil, err
	}

	return &code, nil
}

func (s *Service) Delete(id uint) error {
	return s.db.Delete(&models.Code{}, id).Error
}
