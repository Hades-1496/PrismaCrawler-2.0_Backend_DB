package middlewares

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware configura las cabeceras de origen cruzado
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
		
		if allowedOrigin != "" && origin != allowedOrigin && allowedOrigin != "*" {
			// Si hay origin permitido y no coincide, no ponemos la cabecera (o ponemos la permitida)
			origin = allowedOrigin
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Internal-Token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
