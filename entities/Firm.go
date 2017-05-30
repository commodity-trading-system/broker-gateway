package entities

import "github.com/jinzhu/gorm"

type Firm struct {
	gorm.Model
	Name string
	AccountKey string
	SecretKey string
	Weight int
	CommissionPercent int
}
