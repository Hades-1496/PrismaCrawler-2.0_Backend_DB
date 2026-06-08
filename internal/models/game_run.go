package models

import (
	"time"
)

// GameRun representa una partida individual de un personaje
type GameRun struct {
	ID           uint      `gorm:"primaryKey"`
	CharacterID  uint      `gorm:"not null"`
	Seed         string    // Semilla para la generación del mapa
	CurrentFloor int       `gorm:"default:1"`
	Score        int       `gorm:"default:0"`
	Status       string    `gorm:"default:'In_Progress'"` // In_Progress, Won, Dead
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	// Relaciones
	Character   Character      `gorm:"foreignKey:CharacterID"` // GORM sabrá buscar el personaje
	Inventories []RunInventory `gorm:"foreignKey:RunID"`       // Una run tiene un inventario asociado
}
