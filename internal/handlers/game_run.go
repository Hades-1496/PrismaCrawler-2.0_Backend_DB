package handlers

import (
	"net/http"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"

	"github.com/gin-gonic/gin"
)

// StartRunRequest define los datos necesarios para iniciar la partida
type StartRunRequest struct {
	CharacterID uint `json:"character_id" binding:"required"`
}

// StartRun inicializa una nueva partida para un personaje
func StartRun(c *gin.Context) {
	var req StartRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	userIDValue, _ := c.Get("userID")
	userID := uint(userIDValue.(float64))

	// 1. Verificar que el personaje existe, pertenece al usuario y está vivo
	var character models.Character
	if err := db.DB.Where("id = ? AND user_id = ?", req.CharacterID, userID).First(&character).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personaje no encontrado o no te pertenece"})
		return
	}
	if !character.IsAlive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Este personaje está muerto y no puede iniciar una partida"})
		return
	}

	// 2. Crear la Partida (Run)
	run := models.GameRun{
		CharacterID: character.ID,
		Seed:        utils.GenerateSeed(6), // Genera algo como "X7K2P9"
	}

	if result := db.DB.Create(&run); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al crear la partida"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "¡Partida iniciada!", "run": run})
}

// SaveRunRequest define los datos que Phaser enviará al terminar un piso
type SaveRunRequest struct {
	RunID        uint `json:"run_id" binding:"required"`
	CurrentFloor int  `json:"current_floor" binding:"required"`
	Score        int  `json:"score"`
	CurrentHP    int  `json:"current_hp"`
}

// SaveRun actualiza el estado de la partida y la vida del personaje
func SaveRun(c *gin.Context) {
	var req SaveRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	userIDValue, _ := c.Get("userID")
	userID := uint(userIDValue.(float64))

	var run models.GameRun
	// Preload("Character") es como el 'include: { character: true }' de Prisma
	if err := db.DB.Preload("Character").Where("id = ?", req.RunID).First(&run).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partida no encontrada"})
		return
	}

	// Validaciones de seguridad
	if run.Character.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No tienes permiso para modificar esta partida"})
		return
	}
	if run.Status != "In_Progress" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La partida ya ha finalizado"})
		return
	}

	// Actualizamos los datos
	run.CurrentFloor = req.CurrentFloor
	run.Score = req.Score
	
	// Comprobamos si ha muerto
	if req.CurrentHP <= 0 {
		run.Character.IsAlive = false
		run.Character.BaseHP = 0
		run.Status = "Dead"
	} else {
		run.Character.BaseHP = req.CurrentHP
	}

	// Guardamos ambos modelos en la base de datos
	db.DB.Save(&run)
	db.DB.Save(&run.Character)

	c.JSON(http.StatusOK, gin.H{
		"message": "Progreso guardado correctamente",
		"status":  run.Status,
		"floor":   run.CurrentFloor,
		"hp":      run.Character.BaseHP,
	})
}
