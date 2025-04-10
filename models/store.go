package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Store struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	BrandID     uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"not null"`
	Description string
	IsActive    bool `gorm:"default:true"`
	Metadata    datatypes.JSON
	Rating      float32 `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Brand Brand `gorm:"foreignKey:BrandID;constraint:OnDelete:CASCADE"`
}
