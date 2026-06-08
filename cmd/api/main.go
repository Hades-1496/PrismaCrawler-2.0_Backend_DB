package main

import (
	"log"
	"os"
	"prismacrawler/internal/handlers"
	"prismacrawler/internal/middlewares"
	"prismacrawler/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil { // nil es null para pointers, maps, etc.
		log.Fatal("Error loading .env file")
	}
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8000" // Puerto por defecto si se te olvida ponerlo en el .env
	}
	log.Printf("Iniciando servidor en el puerto %s...", PORT)
	dbURL := os.Getenv("DIRECT_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	db.ConnectDB(dbURL)
	// Routes
	router := gin.Default()

	authGroup := router.Group("/auth")
	authGroup.Use(middlewares.RateLimiter()) // Protegemos las rutas de autenticación
	{
		authGroup.POST("/register", handlers.Register)
		authGroup.POST("/login", handlers.Login)
	}

	// Grupo de rutas del Juego (Protegidas)
	apiGroup := router.Group("/api")
	apiGroup.Use(middlewares.AuthMiddleware()) // Aplicamos el candado a este grupo
	{
		apiGroup.GET("/profile", handlers.GetProfile)
		apiGroup.POST("/characters", handlers.CreateCharacter)
		apiGroup.GET("/characters", handlers.GetCharacters)
		apiGroup.POST("/runs/start", handlers.StartRun)
		apiGroup.PUT("/runs/save", handlers.SaveRun)
		apiGroup.GET("/leaderboard", handlers.GetLeaderboard)
		apiGroup.GET("/items", handlers.GetItems)
	}

	// Grupo de rutas de Administración (Doble protección: Auth + Admin)
	adminGroup := apiGroup.Group("/admin")
	adminGroup.Use(middlewares.AdminMiddleware())
	{
		adminGroup.PUT("/role", handlers.UpdateRole)
	}

	router.Run(":" + PORT)
}
