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
		allowCredentials := false
		switch {
		case rawEnv == "" || rawEnv == "*":
			// Sin allowlist explícita: comodín literal SIN credenciales.
			// (La auth va por header Authorization: Bearer, no por cookies, así que
			//  no se necesitan credenciales cross-origin; y "*" + credentials lo
			//  rechazan los navegadores.)
			responseOrigin = "*"
		case allowedSet[requestOrigin]:
			// Origen explícitamente permitido: lo reflejamos y permitimos credenciales.
			responseOrigin = requestOrigin
			allowCredentials = true
		default:
			// No coincide: devolvemos el primer permitido (el preflight fallará en el navegador).
			for k := range allowedSet {
				responseOrigin = k
				break
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", responseOrigin)
		if allowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Internal-Token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
