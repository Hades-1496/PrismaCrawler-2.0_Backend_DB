package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware configura las cabeceras de origen cruzado.
// ALLOWED_ORIGIN puede ser "*", una URL sin barra final, o varias separadas por coma.
func CORSMiddleware() gin.HandlerFunc {
	rawEnv := os.Getenv("ALLOWED_ORIGIN")

	// Normaliza: quita barras finales y construye el conjunto de orígenes permitidos.
	allowedSet := map[string]bool{}
	for _, o := range strings.Split(rawEnv, ",") {
		o = strings.TrimSpace(strings.TrimRight(o, "/"))
		if o != "" {
			allowedSet[o] = true
		}
	}

	return func(c *gin.Context) {
		requestOrigin := c.GetHeader("Origin")

		var responseOrigin string
		switch {
		case rawEnv == "" || rawEnv == "*":
			// Sin restricción o comodín: reflejamos el origin del request (o *)
			if requestOrigin != "" {
				responseOrigin = requestOrigin
			} else {
				responseOrigin = "*"
			}
		case allowedSet[requestOrigin]:
			// El origin del navegador está en la lista permitida
			responseOrigin = requestOrigin
		default:
			// No coincide: respondemos con el primer origen permitido
			// (el preflight fallará en el navegador, que es el comportamiento correcto)
			for k := range allowedSet {
				responseOrigin = k
				break
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", responseOrigin)
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
