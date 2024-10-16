package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID          uint64 `gorm:"column:id;primary_key;auto_increment;" json:"id"`
	Name        string `gorm:"name;unique;not null"`
	Description string `gorm:"description;not null"`
	Active      bool   `gorm:"active;not null"`
}
