package handlers

import (
	"net/http"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CreateCharacterRequest define el JSON para crear un nuevo héroe
type CreateCharacterRequest struct {
	Name  string `json:"name" binding:"required"`
	Class string `json:"class" binding:"required"` // Ej: "Guerrero", "Mago"
}

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
