package handlers

// --- AUTENTICACIÓN ---
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// --- PERSONAJES ---
type CreateCharacterRequest struct {
	Name  string `json:"name" binding:"required"`
	Class string `json:"class" binding:"required"` // Ej: "Guerrero", "Mago"
}

// --- PARTIDAS (RUNS) ---
type StartRunRequest struct {
	CharacterID uint `json:"character_id" binding:"required"`
	MapID       uint `json:"map_id"` // Opcional, si es 0 se genera procedural
}

type SaveRunRequest struct {
	RunID        uint `json:"run_id" binding:"required"`
	CurrentFloor int  `json:"current_floor" binding:"required"`
	Score        int  `json:"score"`
	CurrentHP    int  `json:"current_hp"`
	Kills        int  `json:"kills"`
	DamageDealt  int  `json:"damage_dealt"`
	DamageTaken  int  `json:"damage_taken"`
}

// --- USUARIOS / ADMIN ---
type UpdateRoleRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=USER ADMIN"`
}
