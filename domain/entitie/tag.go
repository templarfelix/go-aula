package entitie

import (
	"gorm.io/gorm"
	"time"
)

// Tag is representing the Tag data struct
type Tag struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" validate:"required"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
