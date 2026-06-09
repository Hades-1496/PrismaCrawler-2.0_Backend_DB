package handlers

import (
	"net/http"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetEnemies devuelve el catálogo completo de enemigos (bestiario)
func GetEnemies(c *gin.Context) {
	var enemies []models.Enemy
	if err := db.DB.Find(&enemies).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error al obtener el bestiario")
		return
	}
	c.JSON(http.StatusOK, gin.H{"enemies": enemies})
}

// GetMaps devuelve la lista de mapas prefabricados
func GetMaps(c *gin.Context) {
	var maps []models.Map
	if err := db.DB.Find(&maps).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error al obtener la lista de mapas")
		return
	}
	c.JSON(http.StatusOK, gin.H{"maps": maps})
}

// GetMapByID devuelve un mapa prefabricado específico
func GetMapByID(c *gin.Context) {
	var gameMap models.Map
	if err := db.DB.First(&gameMap, c.Param("id")).Error; err != nil {
		utils.SendError(c, http.StatusNotFound, "Mapa no encontrado")
		return
	}
	c.JSON(http.StatusOK, gameMap)
}
