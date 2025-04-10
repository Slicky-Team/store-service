package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type StoreService struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	StoreID      uuid.UUID `gorm:"type:uuid;not null"`
	ServiceName  string    `gorm:"not null"`
	ServicePrice float32   `gorm:"not null"`
	IsActive     bool      `gorm:"default:true"`
	Metadata     datatypes.JSON
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Store Store `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
}
