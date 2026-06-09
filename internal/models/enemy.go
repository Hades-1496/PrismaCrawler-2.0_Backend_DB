package models

type Enemy struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null;uniqueIndex"`
	SpriteKey  string `gorm:"not null"`
	BaseHP     int    `gorm:"default:10"`
	BaseAttack int    `gorm:"default:2"`
	BaseXP     int    `gorm:"default:5"`
}
