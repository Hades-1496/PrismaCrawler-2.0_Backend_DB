package handlers

import (
	"net/http"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetProfile devuelve los datos del usuario y su Top 5 de partidas
func GetProfile(c *gin.Context) {
	// 1. Obtener el ID del usuario desde el JWT
	userIDValue, _ := c.Get("userID")
	userID := uint(userIDValue.(float64))

	// 2. Buscar al usuario
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		utils.SendError(c, http.StatusNotFound, "Usuario no encontrado")
		return
	}

	// 3. Buscar el Top 5 de mejores partidas DE ESTE USUARIO
	// Para ello tenemos que hacer un JOIN con la tabla Characters
	var runs []models.GameRun
	db.DB.Preload("Character").
		Joins("JOIN characters ON characters.id = game_runs.character_id").
		Where("characters.user_id = ? AND game_runs.score > 0", userID).
		Order("game_runs.score desc, game_runs.current_floor desc").
		Limit(5).
		Find(&runs)

	// 4. Mapear las partidas para limpiar la respuesta
	var topRuns []gin.H
	for _, run := range runs {
		topRuns = append(topRuns, gin.H{
			"character": run.Character.Name,
			"class":     run.Character.Class,
			"score":     run.Score,
			"floor":     run.CurrentFloor,
		})
	}

	// 5. Devolver el JSON
	c.JSON(http.StatusOK, gin.H{
		"email":    user.Email,
		"top_runs": topRuns,
	})
}
