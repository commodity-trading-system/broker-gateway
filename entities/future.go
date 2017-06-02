package entities



type Future struct {
	ID int 	`gorm:"primary_key"`
	Name string
	Description string
	Period string
}
