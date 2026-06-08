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
		utils.SendError(c, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	userIDValue, _ := c.Get("userID")
	userID := uint(userIDValue.(float64))

	// 1. Verificar que el personaje existe, pertenece al usuario y está vivo
	var character models.Character
	if err := db.DB.Where("id = ? AND user_id = ?", req.CharacterID, userID).First(&character).Error; err != nil {
		utils.SendError(c, http.StatusNotFound, "Personaje no encontrado o no te pertenece")
		return
	}
	if !character.IsAlive {
		utils.SendError(c, http.StatusBadRequest, "Este personaje está muerto y no puede iniciar una partida")
		return
	}

	// 2. Crear la Partida (Run)
	run := models.GameRun{
		CharacterID: character.ID,
		Seed:        utils.GenerateSeed(6), // Genera algo como "X7K2P9"
	}

	if result := db.DB.Create(&run); result.Error != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error interno al crear la partida")
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
		utils.SendError(c, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	userIDValue, _ := c.Get("userID")
	userID := uint(userIDValue.(float64))

	var run models.GameRun
	// Preload("Character") es como el 'include: { character: true }' de Prisma
	if err := db.DB.Preload("Character").Where("id = ?", req.RunID).First(&run).Error; err != nil {
		utils.SendError(c, http.StatusNotFound, "Partida no encontrada")
		return
	}

	// Validaciones de seguridad
	if run.Character.UserID != userID {
		utils.SendError(c, http.StatusUnauthorized, "No tienes permiso para modificar esta partida")
		return
	}
	if run.Status != "In_Progress" {
		utils.SendError(c, http.StatusBadRequest, "La partida ya ha finalizado")
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

// GetLeaderboard devuelve el Top 10 de mejores partidas globales
func GetLeaderboard(c *gin.Context) {
	var runs []models.GameRun

	// Buscamos el Top 10 ordenado por Puntuación y luego por Piso.
	// Preload("Character") trae los datos del héroe para poder mostrar su nombre.
	db.DB.Preload("Character").
		Where("score > 0"). // Opcional: ignoramos partidas sin puntuar
		Order("score desc, current_floor desc").
		Limit(10).
		Find(&runs)

	// Mapeamos los datos para enviar un JSON limpio a Phaser
	var leaderboard []gin.H
	for _, run := range runs {
		leaderboard = append(leaderboard, gin.H{
			"character": run.Character.Name,
			"class":     run.Character.Class,
			"score":     run.Score,
			"floor":     run.CurrentFloor,
			"status":    run.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{"leaderboard": leaderboard})
}
