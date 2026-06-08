package middlewares

import (
	"net/http"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware verifica que el usuario autenticado tenga el rol ADMIN
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// El AuthMiddleware ya guardó el ID en el contexto, lo recuperamos
		userIDValue, _ := c.Get("userID")
		userID := uint(userIDValue.(float64))

		// Buscamos al usuario en la base de datos para ver su rol
		var user models.User
		if err := db.DB.First(&user, userID).Error; err != nil {
			utils.SendError(c, http.StatusUnauthorized, "Usuario no encontrado")
			c.Abort()
			return
		}

		// Si no es ADMIN, lo bloqueamos
		if user.Role != "ADMIN" {
			utils.SendError(c, http.StatusForbidden, "Acceso denegado: Se requieren permisos de administrador")
			c.Abort()
			return
		}

		c.Next()
	}
}
