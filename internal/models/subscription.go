package models

import (
	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model

	RepositoryID      uint `gorm:"foreignKey:RepositoryID" json:"repository_id"`
	SubscribeCodeID   uint `gorm:"foreignKey:SubscribeCodeID" json:"-"`
	UnsubscribeCodeID uint `json:"-"`

	Email       string `gorm:"type:text;not null;index" json:"email"`
	IsConfirmed bool   `gorm:"type:boolean;not null" json:"confirmed"`
	LastSeenTag string `gorm:"type:text;not null" json:"last_seen_tag"`

	UnsubscribeCode *Code       `gorm:"foreignKey:UnsubscribeCodeID" json:"-"`
	SubscribeCode   *Code       `gorm:"foreignKey:SubscribeCodeID" json:"-"`
	Repository      *Repository `gorm:"foreignKey:RepositoryID" json:"repository,omitempty"`
}
