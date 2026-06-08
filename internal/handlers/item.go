package handlers

import (
	"net/http"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetItems devuelve el catálogo completo de objetos del juego
func GetItems(c *gin.Context) {
	var items []models.Item
	if err := db.DB.Find(&items).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error al obtener el catálogo de objetos")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}
