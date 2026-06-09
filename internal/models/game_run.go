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
	Kills        int       `gorm:"default:0"`             // Enemigos derrotados
	DamageDealt  int       `gorm:"default:0"`             // Daño infligido
	DamageTaken  int       `gorm:"default:0"`             // Daño recibido
	Status       string    `gorm:"default:'In_Progress'"` // In_Progress, Won, Dead
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	// Relaciones
	Character   Character      `gorm:"foreignKey:CharacterID"` // GORM sabrá buscar el personaje
	Inventories []RunInventory `gorm:"foreignKey:RunID"`       // Una run tiene un inventario asociado
}

// UpdateState es un método que encapsula la lógica de juego (SRP), calculando daños y muertes
func (r *GameRun) UpdateState(floor int, score int, currentHP int, kills int, dmgDealt int, dmgTaken int) {
	r.CurrentFloor = floor
	r.Score = score
	r.Kills = kills
	r.DamageDealt = dmgDealt
	r.DamageTaken = dmgTaken

	if currentHP <= 0 {
		r.Character.IsAlive = false
		r.Character.BaseHP = 0
		r.Status = "Dead"
	} else {
		r.Character.BaseHP = currentHP
	}
}
