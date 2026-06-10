package models

// KnowledgeChunk almacena fragmentos de lore o mecánicas para que el chatbot de IA responda (RAG)
type KnowledgeChunk struct {
	ID       uint   `gorm:"primaryKey"`
	Keywords string `gorm:"type:jsonb;default:'[]'"` // Ej: '["goblin", "enemigo", "vida"]'
	Content  string `gorm:"type:text;not null"`      // Ej: "Los goblins son criaturas cobardes con 50 HP..."
	// Podremos añadir el vector 'Embedding' más adelante cuando configuren pgvector
}
