package models

import "time"

// Wallet guarda las monedas y gemas del usuario
type Wallet struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"uniqueIndex;not null" json:"user_id"`
	Coins  int  `gorm:"default:0" json:"coins"`
	Gems   int  `gorm:"default:0" json:"gems"`
}

// WeeklyReward registra cada premio entregado a un operador por quedar en el
// top semanal. El índice único (user_id, year_week) garantiza idempotencia:
// el job de reparto puede ejecutarse varias veces sin pagar dos veces la misma semana.
type WeeklyReward struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_week" json:"user_id"`
	YearWeek  string    `gorm:"not null;uniqueIndex:idx_user_week;size:12" json:"year_week"` // p.ej. "2026-W24"
	Rank      int       `gorm:"not null" json:"rank"`
	Gems      int       `gorm:"not null" json:"gems"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// ProcessedStripeEvent registra cada evento de webhook de Stripe ya procesado.
// El índice único sobre EventID garantiza idempotencia: aunque Stripe reentregue
// el mismo evento, las gemas se acreditan una sola vez.
type ProcessedStripeEvent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EventID   string    `gorm:"uniqueIndex;not null;size:255" json:"event_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Garden guarda el estado del jardín en formato JSON para el frontend
type Garden struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserID uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	Plants string `gorm:"type:jsonb;default:'[]'" json:"plants"`
}
