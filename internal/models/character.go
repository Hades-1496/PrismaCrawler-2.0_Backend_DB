package models

// Character representa la tabla 'characters' (Personajes / Héroes)
type Character struct {
	ID      uint   `gorm:"primaryKey"`
	UserID  uint   `gorm:"not null"` // Llave foránea (Foreign Key) vinculada al User
	User    User   `gorm:"foreignKey:UserID"`
	Name    string `gorm:"not null"`
	Class   string `gorm:"not null"` // Ej: Guerrero, Mago
	Level   int    `gorm:"default:1"`
	BaseHP  int    `gorm:"default:100"`
	IsAlive bool   `gorm:"default:true"`
}
