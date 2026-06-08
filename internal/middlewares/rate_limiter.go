package middlewares

import (
	"net/http"
	"prismacrawler/pkg/utils"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Guardamos un limitador de peticiones por cada dirección IP
var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

// getVisitor devuelve el limitador de una IP o crea uno nuevo si no existe
func getVisitor(ip string) *rate.Limiter {
	mu.Lock()         // Bloqueamos para evitar problemas de concurrencia
	defer mu.Unlock() // Desbloqueamos al salir de la función

	limiter, exists := visitors[ip]
	if !exists {
		// Permitir 1 petición por segundo, con una ráfaga (burst) máxima de 3 peticiones seguidas
		limiter = rate.NewLimiter(1, 3)
		visitors[ip] = limiter
	}
	return limiter
}

// RateLimiter es el middleware que bloquea peticiones si superan el límite
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !getVisitor(c.ClientIP()).Allow() {
			utils.SendError(c, http.StatusTooManyRequests, "Demasiadas peticiones. Por favor, espera un momento.")
			c.Abort()
			return
		}
		c.Next()
	}
}
