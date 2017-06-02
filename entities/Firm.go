package entities


type Firm struct {
	ID int `gorm:"primary_key"`
	Name string
	AccountKey string
	SecretKey string
	Weight int
	CommissionPercent int
}
