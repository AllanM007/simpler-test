package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	// ID          uint    `gorm:"column:id;primary_key;auto_increment;" json:"id"`
	Name        string  `gorm:"name;unique;not null"`
	Description string  `gorm:"description;not null"`
	Price       float64 `gorm:"price;not null"`
	StockLevel  int     `gorm:"stockLevel"`
	Active      bool    `gorm:"active;default:true"`
}
