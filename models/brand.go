package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Brand struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"null"`
	Metadata    datatypes.JSON
	Rating      float32   `gorm:"default:0"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Owner UserProfile `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
}
