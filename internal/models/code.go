package models

import (
	"time"

	"gorm.io/gorm"
)

type CodeType = string

const (
	CodeTypeConfirmation CodeType = "confirmation"
	CodeTypeUnsubscribe  CodeType = "unsubscribe"
)

type Code struct {
	gorm.Model

	SubscriptionID uint `gorm:"foreignKey:SubscriptionID" json:"subscription_id"`

	Code      string    `gorm:"text;not null;unique" json:"code"`
	Type      CodeType  `gorm:"text;not null;index" json:"type"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`

	Subscription *Subscription `gorm:"foreignKey:SubscriptionID" json:"subscription,omitempty"`
}
