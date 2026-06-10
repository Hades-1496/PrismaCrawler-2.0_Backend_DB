package routes

import (
	"prismacrawler/internal/handlers"
	"prismacrawler/internal/middlewares"

	"github.com/gin-gonic/gin"
)

// Setup configura todos los endpoints
func Setup(router *gin.Engine, gameHandler *handlers.GameHandler) {
	router.Use(middlewares.CORSMiddleware())

	// Leaderboard PÚBLICO (sin JWT)
	router.GET("/api/leaderboard", gameHandler.GetLeaderboard)

	// Grupo Auth
	authGroup := router.Group("/auth")
	authGroup.Use(middlewares.RateLimiter())
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
		apiGroup.POST("/runs/start", gameHandler.StartRun)
		apiGroup.PUT("/runs/save", gameHandler.SaveRun)
		apiGroup.GET("/items", handlers.GetItems)
		apiGroup.GET("/enemies", handlers.GetEnemies)
		apiGroup.GET("/maps", handlers.GetMaps)
		apiGroup.GET("/maps/:id", handlers.GetMapByID)
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
}
