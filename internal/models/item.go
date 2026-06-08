package models

// Item representa el catálogo global de objetos en el juego
type Item struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	SpriteKey   string // Para mantener compatibilidad con tu frontend (ej: 'item_sword')
	Description string
	Type        string `gorm:"not null"` // Ej: "Weapon", "Armor", "Potion"
	Rarity      string `gorm:"default:'common'"`

	// JSONB es súper potente en PostgreSQL. Nos permite guardar un objeto JSON flexible
	StatsModifier string `gorm:"type:jsonb"`
	IsConsumable  bool   `gorm:"default:false"`
}
