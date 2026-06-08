package models

// RunInventory actúa como tabla intermedia/detalle entre GameRun e Item
type RunInventory struct {
	ID       uint `gorm:"primaryKey"`
	RunID    uint `gorm:"not null"`
	ItemID   uint `gorm:"not null"`
	Quantity int  `gorm:"default:1"`
	Equipped bool `gorm:"default:false"`

	// Relación
	// Nos permitirá hacer un ".Preload('Item')" más adelante para traer los datos del objeto
	Item Item `gorm:"foreignKey:ItemID"`
}
