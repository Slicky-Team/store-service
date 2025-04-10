package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	Admin   UserRole = "ADMIN"
	Manager UserRole = "MANAGER"
	Staff   UserRole = "STAFF"
)

type UserStore struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	StoreID   uuid.UUID `gorm:"type:uuid;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"` // references UserProfile.AccountID
	Role      UserRole  `gorm:"type:user_role;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Store Store       `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	User  UserProfile `gorm:"foreignKey:AccountID;references:UserID;constraint:OnDelete:CASCADE"`
}
