package middlewares

import "github.com/gin-gonic/gin"

// CORSMiddleware permite que un frontend en otro dominio/puerto haga peticiones a esta API
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// En producción, podrías cambiar "*" por la URL exacta de tu frontend en Vercel/Netlify
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// Interceptamos las peticiones de "Preflight" (OPTIONS) que hacen los navegadores
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
