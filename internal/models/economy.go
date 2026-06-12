package models

// Wallet guarda las monedas y gemas del usuario
type Wallet struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"uniqueIndex;not null" json:"user_id"`
	Coins  int  `gorm:"default:0" json:"coins"`
	Gems   int  `gorm:"default:0" json:"gems"`
}

// Garden guarda el estado del jardín en formato JSON para el frontend
type Garden struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserID uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	Plants string `gorm:"type:jsonb;default:'[]'" json:"plants"`
}
