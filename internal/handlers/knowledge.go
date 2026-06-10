package handlers

import (
	"encoding/json"
	"net/http"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CreateKnowledge añade un nuevo fragmento de lore a la base de datos
func CreateKnowledge(c *gin.Context) {
	var req KnowledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// Convertimos el array de strings a un JSON stringificado para guardarlo en la columna JSONB
	kwBytes, _ := json.Marshal(req.Keywords)

	chunk := models.KnowledgeChunk{
		Keywords: string(kwBytes),
		Content:  req.Content,
	}

	if err := db.DB.Create(&chunk).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error al guardar el conocimiento")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Conocimiento guardado exitosamente", "knowledge": chunk})
}

// GetKnowledge devuelve todo el conocimiento guardado para revisar el lore
func GetKnowledge(c *gin.Context) {
	var chunks []models.KnowledgeChunk
	db.DB.Find(&chunks)
	c.JSON(http.StatusOK, gin.H{"knowledge": chunks})
}

// DeleteKnowledge elimina un fragmento por su ID
func DeleteKnowledge(c *gin.Context) {
	id := c.Param("id")
	db.DB.Delete(&models.KnowledgeChunk{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Fragmento eliminado correctamente"})
}
