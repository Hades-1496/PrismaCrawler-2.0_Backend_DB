package middlewares

import (
	"net/http"
	"prismacrawler/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware protege las rutas verificando el JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Buscamos el token en la cabecera "Authorization"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Acceso denegado. No hay token."})
			c.Abort() // ¡Súper importante en Gin para detener la cadena de ejecución!
			return
		}

		// 2. Limpiamos el "Bearer " de la cadena
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// 3. Verificamos el token usando nuestro paquete utils
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token no válido"})
			c.Abort()
			return
		}

		// 4. Guardamos los datos en el contexto (equivalente a req.user = verified)
		c.Set("userID", claims["user_id"])

		// 5. NEXT(): Permite que la petición continúe hacia el controlador (ej. CreateCharacter)
		c.Next()
	}
}
