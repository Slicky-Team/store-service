package models

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	AccountID   uuid.UUID `gorm:"type:uuid;unique;not null"`
	FullName    string
	Age         int
	PhoneNumber string
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
