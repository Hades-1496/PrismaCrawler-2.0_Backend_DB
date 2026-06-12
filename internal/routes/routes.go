package routes

import (
	"net/http"
	"os"
	"prismacrawler/internal/handlers"
	"prismacrawler/internal/middlewares"
	"time"

	"github.com/gin-gonic/gin"
)

// Setup configura todos los endpoints
func Setup(router *gin.Engine, gameHandler *handlers.GameHandler, internalHandler *handlers.InternalHandler) {
	router.Use(middlewares.CORSMiddleware())

	// Leaderboard PÚBLICO (sin JWT)
	router.GET("/api/leaderboard", gameHandler.GetLeaderboard)

	// Health check (Ping) PÚBLICO para comprobar que el deploy está funcionando
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Grupo Auth
	authGroup := router.Group("/auth")
	authGroup.Use(middlewares.RateLimiter())
	authGroup.Use(func(c *gin.Context) {
		// Retraso artificial para mitigar ataques de timing y fuerza bruta
		time.Sleep(300 * time.Millisecond)
		c.Next()
	})
	{
		authGroup.POST("/register", handlers.Register)
		authGroup.POST("/login", handlers.Login)
	}

	// Grupo API (Core protegido por JWT)
	apiGroup := router.Group("/api")
	apiGroup.Use(middlewares.AuthMiddleware())
	{
		apiGroup.GET("/profile", handlers.GetProfile)
		apiGroup.POST("/characters", handlers.CreateCharacter)
		apiGroup.GET("/characters", handlers.GetCharacters)
		apiGroup.GET("/characters/:id", handlers.GetCharacterByID)
		apiGroup.PUT("/characters/:id", handlers.UpdateCharacter)
		apiGroup.DELETE("/characters/:id", handlers.DeleteCharacter)
		apiGroup.POST("/runs/start", gameHandler.StartRun)
		apiGroup.PUT("/runs/save", gameHandler.SaveRun)
		apiGroup.GET("/runs", gameHandler.GetRuns)
		apiGroup.GET("/runs/:id", gameHandler.GetRunByID)
		apiGroup.GET("/items", handlers.GetItems)
		apiGroup.GET("/enemies", handlers.GetEnemies)
		apiGroup.GET("/maps", handlers.GetMaps)
		apiGroup.GET("/maps/:id", handlers.GetMapByID)
		apiGroup.GET("/wallet", handlers.GetWallet)
		apiGroup.PUT("/wallet", handlers.UpdateWallet)
		apiGroup.GET("/garden", handlers.GetGarden)
		apiGroup.PUT("/garden", handlers.UpdateGarden)
	}

	// Grupo Admin (Doble protección: Auth + Admin)
	adminGroup := apiGroup.Group("/admin")
	adminGroup.Use(middlewares.AdminMiddleware())
	{
		adminGroup.PUT("/role", handlers.UpdateRole)
		adminGroup.POST("/discord/changelog", handlers.DiscordChangelog)
		adminGroup.POST("/discord/test-webhook", handlers.DiscordTestWebhook)
		adminGroup.POST("/knowledge", handlers.CreateKnowledge)
		adminGroup.GET("/knowledge", handlers.GetKnowledge)
		adminGroup.DELETE("/knowledge/:id", handlers.DeleteKnowledge)
	}

	// Alias para mitigar el Bug 2 del frontend (intentaban consumir /game/items)
	router.GET("/game/items", handlers.GetItems)

	// Grupo de rutas internas exclusivas para el microservicio IA
	internalGroup := router.Group("/api/internal")
	internalGroup.Use(func(c *gin.Context) {
		if c.GetHeader("X-Internal-Token") != os.Getenv("AI_INTERNAL_TOKEN") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Acceso denegado"})
			return
		}
		c.Next()
	})
	{
		internalGroup.GET("/rankings/top10", internalHandler.GetTop10)
		internalGroup.POST("/rag/search", internalHandler.SearchRAG)
		internalGroup.GET("/stats", internalHandler.GetStats)
		internalGroup.GET("/items", internalHandler.GetItems)
		internalGroup.GET("/enemies", internalHandler.GetEnemies)
		internalGroup.GET("/wallet/:user_id", internalHandler.GetWallet)
		internalGroup.GET("/garden/:user_id", internalHandler.GetGarden)
	}
}
