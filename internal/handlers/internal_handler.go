package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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

	result := db.DB.Preload("Character.User").
		Where("status IN ?", []string{"Dead", "Won"}).
		Order("score DESC, current_floor DESC").
		Limit(10).
		Find(&topRuns)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el ranking"})
		return
	}

	var leaderboard []gin.H
	for _, run := range topRuns {
		email := run.Character.User.Email
		username := strings.SplitN(email, "@", 2)[0]
		leaderboard = append(leaderboard, gin.H{
			"character_id":   run.CharacterID,
			"character_name": run.Character.Name,
			"username":       username,
			"user_id":        run.Character.UserID,
			"score":          run.Score,
			"floor":          run.CurrentFloor,
			"kills":          run.Kills,
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

// gemsForRank define el premio escalonado en GEMAS (moneda premium) por posición
// en el top semanal. Las gemas no se generan en el garden ni en la mazmorra, así
// que repartirlas como premio de ranking NO infla la economía de monedas (coins).
func gemsForRank(rank int) int {
	switch {
	case rank == 1:
		return 50
	case rank == 2:
		return 20
	case rank == 3:
		return 15
	default:
		return 10
	}
}

// mondayOfWeekUTC devuelve el lunes a las 00:00 UTC de la semana que contiene t.
func mondayOfWeekUTC(t time.Time) time.Time {
	t = t.UTC()
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // domingo cuenta como último día de la semana ISO
	}
	monday := t.AddDate(0, 0, -(weekday - 1))
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)
}

// DistributeWeeklyRewards reparte gemas al top 10 de operadores de una semana.
// Por defecto procesa la SEMANA ANTERIOR (cierre de competición). Con ?week=current
// procesa la semana en curso (útil para pruebas/demostración).
// Es idempotente: el índice único (user_id, year_week) impide pagar dos veces.
func (h *InternalHandler) DistributeWeeklyRewards(c *gin.Context) {
	now := time.Now().UTC()
	target := now.AddDate(0, 0, -7) // semana anterior por defecto
	if c.Query("week") == "current" {
		target = now
	}

	year, week := target.ISOWeek()
	yearWeek := fmt.Sprintf("%d-W%02d", year, week)
	weekStart := mondayOfWeekUTC(target)
	weekEnd := weekStart.AddDate(0, 0, 7)

	// Mejores runs de la semana. Pedimos margen (50) para luego deduplicar por usuario:
	// un mismo operador no puede ocupar varias posiciones del podio.
	var topRuns []models.GameRun
	if err := db.DB.Preload("Character.User").
		Where("created_at >= ? AND created_at < ?", weekStart, weekEnd).
		Order("score DESC, current_floor DESC").
		Limit(50).
		Find(&topRuns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar el ranking semanal"})
		return
	}

	seen := map[uint]bool{}
	rank := 0
	awarded := []gin.H{}

	for _, run := range topRuns {
		uid := run.Character.UserID
		if uid == 0 || seen[uid] {
			continue
		}
		seen[uid] = true
		rank++
		if rank > 10 {
			break
		}

		gems := gemsForRank(rank)
		reward := models.WeeklyReward{UserID: uid, YearWeek: yearWeek, Rank: rank, Gems: gems}

		// El Create falla si ya existe (user_id, year_week) → semana ya pagada, lo saltamos.
		if err := db.DB.Create(&reward).Error; err != nil {
			continue
		}

		// Acreditamos las gemas en la billetera del operador.
		var wallet models.Wallet
		db.DB.Where("user_id = ?", uid).FirstOrCreate(&wallet, models.Wallet{UserID: uid})
		wallet.Gems += gems
		db.DB.Save(&wallet)

		nickname := run.Character.User.Nickname
		if nickname == "" {
			nickname = run.Character.Name
		}
		awarded = append(awarded, gin.H{
			"user_id":        uid,
			"nickname":       nickname,
			"character_name": run.Character.Name,
			"rank":           rank,
			"gems":           gems,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"year_week":     yearWeek,
		"week_start":    weekStart,
		"week_end":      weekEnd,
		"awarded_count": len(awarded),
		"awarded":       awarded,
	})
}
