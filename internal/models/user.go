package models

import (
	"time"
)

// User representa la tabla 'users' en PostgreSQL
type User struct {
	ID           uint        `gorm:"primaryKey"`
	Email        string      `gorm:"uniqueIndex;not null"`
	PasswordHash string      `gorm:"not null"`
	Role         string      `gorm:"default:'USER'"` // Rol del usuario (ej: USER, ADMIN)
	Nickname     string      `gorm:"uniqueIndex;not null;size:32"`
	RealName     string      `gorm:"size:64"`
	Avatar       string      `gorm:"size:128"`
	PlayerIcon   string      `gorm:"size:64"`
	CreatedAt    time.Time   `gorm:"autoCreateTime"`
	LastLogin    time.Time   // Sin etiqueta especial, GORM asume que es una columna normal (timestamp)
	Characters   []Character `gorm:"foreignKey:UserID"` // Relación: Un User tiene muchos Characters
}
