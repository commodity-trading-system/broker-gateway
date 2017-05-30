package entities

import (
	"github.com/jinzhu/gorm"
)

type Future struct {
	gorm.Model
	Name string
	Description string
	Period string
}
