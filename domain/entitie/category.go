package entitie

import (
	"gorm.io/gorm"
)

type Category struct {
	Name string `json:"name" validate:"required"`
	gorm.Model
}
