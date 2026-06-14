package handlers

// --- AUTENTICACIÓN ---
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname" binding:"required,min=3,max=32"`
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
	CharacterID uint `json:"character_id"` // Quitamos binding:"required" para mitigar bug de Phaser
	MapID       uint `json:"map_id"`       // Opcional, si es 0 se genera procedural
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

// --- IA / CHATBOT ---
type FaqRequest struct {
	Question string `json:"question" binding:"required,min=3,max=500"`
}

// --- IA / DISCORD (admin) ---
type ChangelogRequest struct {
	Title   string `json:"title" binding:"max=256"`
	Content string `json:"content" binding:"required,min=1,max=4000"`
}

// --- USUARIOS / ADMIN ---
type UpdateRoleRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=USER ADMIN"`
}

type UpdateProfileRequest struct {
	Nickname   string `json:"nickname" binding:"max=32"`
	RealName   string `json:"real_name" binding:"max=64"`
	Avatar     string `json:"avatar" binding:"max=128"`
	PlayerIcon string `json:"player_icon" binding:"max=64"`
	Role       string `json:"role" binding:"omitempty,oneof=USER ADMIN"`
}

// --- ECONOMÍA ---
type UpdateWalletRequest struct {
	DeltaCoins int `json:"delta_coins"`
	DeltaGems  int `json:"delta_gems"`
}

type UpdateGardenRequest struct {
	Plants string `json:"plants" binding:"required"`
}

type CreateCheckoutRequest struct {
	PackageID string `json:"package_id" binding:"required"` // e.g. "gems_500"
}

type ExchangeCoinsRequest struct {
	CoinsToSpend int `json:"coins_to_spend" binding:"required,min=1000"` // debe ser múltiplo de 1000
}

// --- CONOCIMIENTO (RAG / ADMIN) ---
type KnowledgeRequest struct {
	Keywords []string `json:"keywords" binding:"required"`
	Content  string   `json:"content" binding:"required,min=10"`
}

// --- IA / RAG ---
type RAGSearchRequest struct {
	Embedding []float32 `json:"embedding" binding:"required"`
}
