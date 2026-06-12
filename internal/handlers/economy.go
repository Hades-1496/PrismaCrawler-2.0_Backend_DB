package handlers

import (
	"context"
	"net/http"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// GetWallet devuelve la billetera del usuario logueado
func GetWallet(c *gin.Context) {
	userID := utils.GetUserID(c)
	var wallet models.Wallet
	db.DB.Where("user_id = ?", userID).FirstOrCreate(&wallet, models.Wallet{UserID: userID})
	c.JSON(http.StatusOK, wallet)
}

// UpdateWallet actualiza las monedas y gemas del usuario logueado
func UpdateWallet(c *gin.Context) {
	userID := utils.GetUserID(c)
	var req models.Wallet
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Datos de billetera inválidos")
		return
	}

	var wallet models.Wallet
	db.DB.Where("user_id = ?", userID).FirstOrCreate(&wallet, models.Wallet{UserID: userID})
	wallet.Coins = req.Coins
	wallet.Gems = req.Gems
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

// UpdateGarden actualiza las plantas del usuario logueado
func UpdateGarden(c *gin.Context) {
	userID := utils.GetUserID(c)
	var req models.Garden
	var garden models.Garden
	db.DB.Where("user_id = ?", userID).FirstOrCreate(&garden, models.Garden{UserID: userID})
	c.ShouldBindJSON(&req) // Si falla usaremos los valores de req por defecto (vacíos)
	garden.Plants = req.Plants
	db.DB.Save(&garden)
	c.JSON(http.StatusOK, garden)
}
