package entitie

import (
	"gorm.io/gorm"
)

type Tag struct {
	Name string `json:"name" validate:"required"`
	gorm.Model
}
