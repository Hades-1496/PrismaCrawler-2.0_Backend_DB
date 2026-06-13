package handlers

import (
	"net/http"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type InternalHandler struct{}

func NewInternalHandler() *InternalHandler {
	return &InternalHandler{}
}

func (h *InternalHandler) GetTop10(c *gin.Context) {
	var topRuns []models.GameRun

	// Obtenemos las 10 mejores partidas (ordenadas por score y piso)
	result := db.DB.Preload("Character").
		Order("score DESC, current_floor DESC").
		Limit(10).
		Find(&topRuns)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el ranking"})
		return
	}

	var leaderboard []gin.H
	for _, run := range topRuns {
		leaderboard = append(leaderboard, gin.H{
			"character_id":   run.CharacterID,
			"character_name": run.Character.Name,
			"user_id":        run.Character.UserID,
			"score":          run.Score,
			"floor":          run.CurrentFloor,
		})
	}

	c.JSON(http.StatusOK, leaderboard)
}

// SearchRAG recibe un vector desde Python y devuelve los fragmentos más relevantes usando pgvector
func (h *InternalHandler) SearchRAG(c *gin.Context) {
	var req RAGSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vector embedding requerido"})
		return
	}

	var chunks []models.KnowledgeChunk
	// pgvector: ordenamos por similitud (distancia L2 usando el operador <->)
	// Limitamos a los 3 mejores resultados para inyectar como contexto al chatbot.
	result := db.DB.Order(gorm.Expr("embedding <-> ?", pgvector.NewVector(req.Embedding))).Limit(3).Find(&chunks)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en búsqueda vectorial"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": chunks})
}

// GetStats devuelve estadísticas globales del juego
func (h *InternalHandler) GetStats(c *gin.Context) {
	var totalPlayers, totalRuns, totalDeaths int64

	db.DB.Model(&models.User{}).Count(&totalPlayers)
	db.DB.Model(&models.GameRun{}).Count(&totalRuns)
	db.DB.Model(&models.GameRun{}).Where("status = ?", "Dead").Count(&totalDeaths)

	c.JSON(http.StatusOK, gin.H{
		"total_players": totalPlayers,
		"total_runs":    totalRuns,
		"total_deaths":  totalDeaths,
	})
}

// GetItems expone el catálogo exclusivo para la IA
func (h *InternalHandler) GetItems(c *gin.Context) {
	var items []models.Item
	db.DB.Find(&items)
	c.JSON(http.StatusOK, items)
}

// GetEnemies expone el bestiario exclusivo para la IA
func (h *InternalHandler) GetEnemies(c *gin.Context) {
	var enemies []models.Enemy
	db.DB.Find(&enemies)
	c.JSON(http.StatusOK, enemies)
}

// GetWallet devuelve la economía de un usuario específico para la IA
func (h *InternalHandler) GetWallet(c *gin.Context) {
	userID := c.Param("user_id")
	var wallet models.Wallet
	if err := db.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"coins": 0, "gems": 0})
		return
	}
	c.JSON(http.StatusOK, wallet)
}

// GetGarden devuelve el jardín de un usuario específico para la IA
func (h *InternalHandler) GetGarden(c *gin.Context) {
	userID := c.Param("user_id")
	var garden models.Garden
	if err := db.DB.Where("user_id = ?", userID).First(&garden).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"plants": "[]"})
		return
	}
	c.JSON(http.StatusOK, garden)
}
