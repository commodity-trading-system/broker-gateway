package entities

type Commission struct {
	ID int	`gorm:"primary_key"`
	FirmId int
	OrderType int
	FutureId int
	CommissionPercent int
}
