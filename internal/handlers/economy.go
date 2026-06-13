package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	maxDeltaCoinsPerTx = 50_000
	maxDeltaGemsPerTx  = 5_000
	maxCoinsTotal      = 10_000_000
	maxGemsTotal       = 1_000_000
)

// GetWallet devuelve la billetera del usuario logueado
func GetWallet(c *gin.Context) {
	userID := utils.GetUserID(c)
	var wallet models.Wallet
	db.DB.Where("user_id = ?", userID).FirstOrCreate(&wallet, models.Wallet{UserID: userID})
	c.JSON(http.StatusOK, wallet)
}

// UpdateWallet aplica un delta de monedas/gemas en lugar de sobrescribir el total.
// Previene mass-assignment: el cliente no puede establecer un valor arbitrario directamente.
func UpdateWallet(c *gin.Context) {
	userID := utils.GetUserID(c)
	var req UpdateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Datos de billetera inválidos")
		return
	}

	if req.DeltaCoins > maxDeltaCoinsPerTx || req.DeltaGems > maxDeltaGemsPerTx {
		utils.SendError(c, http.StatusBadRequest, "Delta fuera del rango permitido")
		return
	}

	var wallet models.Wallet
	db.DB.Where("user_id = ?", userID).FirstOrCreate(&wallet, models.Wallet{UserID: userID})

	wallet.Coins = max(0, min(maxCoinsTotal, wallet.Coins+req.DeltaCoins))
	wallet.Gems = max(0, min(maxGemsTotal, wallet.Gems+req.DeltaGems))
	db.DB.Save(&wallet)

	// Notificamos asíncronamente al módulo de Economía de la IA
	if AI != nil {
		go func(uID uint, coins, gems int) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			payload := map[string]any{
				"event": "wallet_updated",
				"data": map[string]any{
					"user_id": uID,
					"coins":   coins,
					"gems":    gems,
				},
			}
			AI.Proxy(ctx, "/api/economy/status", payload)
		}(userID, wallet.Coins, wallet.Gems)
	}

	c.JSON(http.StatusOK, wallet)
}

// GetGarden devuelve el jardín del usuario logueado
func GetGarden(c *gin.Context) {
	userID := utils.GetUserID(c)
	var garden models.Garden
	db.DB.Where("user_id = ?", userID).FirstOrCreate(&garden, models.Garden{UserID: userID})
	c.JSON(http.StatusOK, garden)
}

// UpdateGarden actualiza las plantas del usuario logueado.
// Valida que el campo plants sea JSON válido antes de persistir.
func UpdateGarden(c *gin.Context) {
	userID := utils.GetUserID(c)
	var req UpdateGardenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Datos de jardín inválidos")
		return
	}
	if !json.Valid([]byte(req.Plants)) {
		utils.SendError(c, http.StatusBadRequest, "El campo plants no contiene JSON válido")
		return
	}

	var garden models.Garden
	db.DB.Where("user_id = ?", userID).FirstOrCreate(&garden, models.Garden{UserID: userID})
	garden.Plants = req.Plants
	db.DB.Save(&garden)
	c.JSON(http.StatusOK, garden)
}
