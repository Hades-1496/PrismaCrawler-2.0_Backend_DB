package handlers

import (
	"net/http"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CreateCharacter maneja la creación de un personaje asociado al usuario logueado
func CreateCharacter(c *gin.Context) {
	var req CreateCharacterRequest

	// Validamos el JSON recibido
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// Extraemos el ID del usuario desde el contexto (colocado allí por el AuthMiddleware)
	userID := utils.GetUserID(c)

	character := models.Character{
		UserID: userID,
		Name:   req.Name,
		Class:  req.Class,
		// Level, BaseHP e IsAlive tomarán los valores por defecto de la base de datos
	}

	if result := db.DB.Create(&character); result.Error != nil {
		utils.SendError(c, http.StatusInternalServerError, "No se pudo crear el personaje")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Personaje creado con éxito", "character": character})
}

// GetCharacters devuelve la lista de personajes del usuario logueado
func GetCharacters(c *gin.Context) {
	userID := utils.GetUserID(c)

	var characters []models.Character
	// Buscamos todos los personajes que pertenezcan a este usuario
	db.DB.Where("user_id = ?", userID).Find(&characters)

	// Si no hay personajes, devolverá un array vacío [], lo cual es correcto
	c.JSON(http.StatusOK, gin.H{"characters": characters})
}

// GetCharacterByID devuelve un personaje específico del usuario
func GetCharacterByID(c *gin.Context) {
	userID := utils.GetUserID(c)
	var character models.Character

	if err := db.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&character).Error; err != nil {
		utils.SendError(c, http.StatusNotFound, "Personaje no encontrado")
		return
	}

	c.JSON(http.StatusOK, character)
}

// UpdateCharacter actualiza la información de un personaje
func UpdateCharacter(c *gin.Context) {
	userID := utils.GetUserID(c)
	var req CreateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	var character models.Character
	if err := db.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&character).Error; err != nil {
		utils.SendError(c, http.StatusNotFound, "Personaje no encontrado")
		return
	}

	character.Name = req.Name
	character.Class = req.Class
	db.DB.Save(&character)

	c.JSON(http.StatusOK, gin.H{"message": "Personaje actualizado", "character": character})
}

// DeleteCharacter elimina un personaje si pertenece al usuario
func DeleteCharacter(c *gin.Context) {
	userID := utils.GetUserID(c)
	result := db.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).Delete(&models.Character{})
	if result.Error != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error al intentar eliminar personaje")
		return
	}
	if result.RowsAffected == 0 {
		utils.SendError(c, http.StatusNotFound, "Personaje no encontrado o no autorizado")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Personaje eliminado correctamente"})
}
