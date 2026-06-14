package handlers

import (
	"net/http"
	"prismacrawler/internal/models"
	"prismacrawler/internal/services"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"
	"strconv"

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

	// --- MITIGACIÓN BUG FRONTEND ---
	// Si el frontend envía body: {} sin character_id, auto-seleccionamos el primer personaje
	if req.CharacterID == 0 {
		var firstChar models.Character
		db.DB.Where("user_id = ?", userID).First(&firstChar)
		if firstChar.ID != 0 {
			req.CharacterID = firstChar.ID
		} else {
			// Auto-crear personaje por defecto si no tiene ninguno
			firstChar = models.Character{
				UserID: userID,
				Name:   "Operador",
				Class:  "Soldier",
			}
			db.DB.Create(&firstChar)
			req.CharacterID = firstChar.ID
		}
	}

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
		playerName := run.Character.Name
		if run.Character.User.Nickname != "" {
			playerName = run.Character.User.Nickname
		}
		leaderboard = append(leaderboard, gin.H{
			"playerName":       playerName,
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

// GetRuns devuelve el historial de todas las partidas del usuario (paginado)
func (h *GameHandler) GetRuns(c *gin.Context) {
	userID := utils.GetUserID(c)

	limitVal := 10
	offsetVal := 0

	if limitParam := c.Query("limit"); limitParam != "" {
		if val, err := strconv.Atoi(limitParam); err == nil && val > 0 {
			limitVal = val
		}
	}
	if pageParam := c.Query("page"); pageParam != "" {
		if val, err := strconv.Atoi(pageParam); err == nil && val > 0 {
			offsetVal = (val - 1) * limitVal
		}
	}

	var runs []models.GameRun

	// Hacemos JOIN para asegurar que la partida le pertenece al usuario a través del personaje
	db.DB.Preload("Character").
		Joins("JOIN characters ON characters.id = game_runs.character_id").
		Where("characters.user_id = ?", userID).
		Order("game_runs.created_at desc").
		Limit(limitVal).
		Offset(offsetVal).
		Find(&runs)

	c.JSON(http.StatusOK, gin.H{"runs": runs})
}

// GetRunByID devuelve los detalles de una partida, incluyendo el inventario (RunInventory)
func (h *GameHandler) GetRunByID(c *gin.Context) {
	userID := utils.GetUserID(c)
	var run models.GameRun

	// Anidamos el Preload para obtener el inventario y los detalles de los objetos en un solo JSON
	if err := db.DB.Preload("Character").Preload("Inventories.Item").Where("id = ?", c.Param("id")).First(&run).Error; err != nil {
		utils.SendError(c, http.StatusNotFound, "Partida no encontrada")
		return
	}

	if !run.IsOwnedBy(userID) {
		utils.SendError(c, http.StatusForbidden, "No tienes permiso para ver esta partida")
		return
	}

	c.JSON(http.StatusOK, run)
}

// GetEconomyLeaderboard devuelve el top 10 de operadores más ricos
func (h *GameHandler) GetEconomyLeaderboard(c *gin.Context) {
	var wallets []models.Wallet
	err := db.DB.Order("coins desc, gems desc").Limit(10).Find(&wallets).Error
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error al obtener el ranking de economía")
		return
	}

	var leaderboard []gin.H
	for _, wallet := range wallets {
		var user models.User
		db.DB.First(&user, wallet.UserID)
		
		playerName := "Operador"
		if user.Nickname != "" {
			playerName = user.Nickname
		} else if user.Email != "" {
			playerName = user.Email
		}

		leaderboard = append(leaderboard, gin.H{
			"playerName": playerName,
			"coins":      wallet.Coins,
			"gems":       wallet.Gems,
		})
	}

	c.JSON(http.StatusOK, gin.H{"leaderboard": leaderboard})
}
