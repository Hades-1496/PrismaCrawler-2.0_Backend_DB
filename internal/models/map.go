package models

type Map struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	Layout     string `gorm:"type:jsonb"` // El layout ASCII como JSON
	Dictionary string `gorm:"type:jsonb"` // El diccionario de entidades como JSON
}
