package utils

import (
	"github.com/gin-gonic/gin"
)

// SendError estandariza la forma en la que la API devuelve errores
func SendError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
