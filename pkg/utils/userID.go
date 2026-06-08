package utils

import "github.com/gin-gonic/gin"

// GetUserID extrae de forma segura el ID del usuario del contexto (inyectado por el AuthMiddleware)
func GetUserID(c *gin.Context) uint {
	userIDValue, exists := c.Get("userID")
	if !exists {
		return 0
	}
	// El JWT parsea los números como float64, los convertimos a uint
	if id, ok := userIDValue.(float64); ok {
		return uint(id)
	}
	return 0
}
