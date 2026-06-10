package handlers

import (
	
	"net/http"
	"prismacrawler/internal/services"
	"prismacrawler/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GameHandler agrupa los controladores de partidas y sus dependencias (el servicio).
type GameHandler struct {
	service services.GameServiceInterface
}

func NewGameHandler(service services.GameServiceInterface) *GameHandler {
	return &GameHandler{service: service}
}

// StartRun inicializa una nueva partida para un personaje
func (h *GameHandler) StartRun(c *gin.Context) {
	var req StartRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	userID := utils.GetUserID(c)
	input := services.StartRunInput{
		CharacterID: req.CharacterID,
		MapID:       req.MapID,
	}

	run, err := h.service.StartRun(c.Request.Context(), userID, input)
	if err != nil {
		switch err {
		case services.ErrCharacterNotFound:
			utils.SendError(c, http.StatusNotFound, err.Error())
		case services.ErrCharacterDead:
			utils.SendError(c, http.StatusBadRequest, err.Error())
		default:
			utils.SendError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "¡Partida iniciada!", "run": run})
}

// SaveRun actualiza el estado de la partida y la vida del personaje
func (h *GameHandler) SaveRun(c *gin.Context) {
	var req SaveRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}
	userID := utils.GetUserID(c)

	// Mapeamos el DTO de HTTP al DTO del Servicio
	input := services.SaveRunInput{
		RunID:        req.RunID,
		CurrentFloor: req.CurrentFloor,
		Score:        req.Score,
		CurrentHP:    req.CurrentHP,
		Kills:        req.Kills,
		DamageDealt:  req.DamageDealt,
		DamageTaken:  req.DamageTaken,
	}

	run, err := h.service.SaveRun(c.Request.Context(), userID, input)
	if err != nil {
		switch err {
		case services.ErrRunNotFound:
			utils.SendError(c, http.StatusNotFound, err.Error())
		case services.ErrPermissionDenied:
			utils.SendError(c, http.StatusForbidden, err.Error())
		case services.ErrRunAlreadyEnded:
			utils.SendError(c, http.StatusBadRequest, err.Error())
		default:
			utils.SendError(c, http.StatusInternalServerError, "Error al guardar la partida")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Progreso guardado correctamente",
		"status":  run.Status,
		"floor":   run.CurrentFloor,
		"hp":      run.Character.BaseHP,
	})
}

// GetLeaderboard devuelve el Top 10 de mejores partidas globales
func (h *GameHandler) GetLeaderboard(c *gin.Context) {
	runs, err := h.service.GetLeaderboard(c.Request.Context())
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error al obtener el leaderboard")
		return
	}

	// Mapeamos los datos para enviar un JSON limpio a Phaser
	var leaderboard []gin.H
	for _, run := range runs {
		leaderboard = append(leaderboard, gin.H{
			"playerName":       run.Character.Name,
			"class":            run.Character.Class,
			"score":            run.Score,
			"floor":            run.CurrentFloor,
			"kills":            run.Kills,
			"totalDamageDealt": run.DamageDealt,
			"totalDamageTaken": run.DamageTaken,
			"status":           run.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{"leaderboard": leaderboard})
}
