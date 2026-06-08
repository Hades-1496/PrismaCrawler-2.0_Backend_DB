package main

import (
	"log" // Mirar diferencias entre log y fmt
	"net/http"
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
	db.ConnectDB(os.Getenv("DATABASE_URL"))
	// Routes
	router := gin.Default()

	router.GET("/ping", getting)

	authGroup := router.Group("/auth")
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
	}

	router.Run(":" + PORT)
}

// Controller?
func getting(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"method": "GET"})
}

// Objetivos adicionales: Dividir el archivo en varios como routes, controllers, utils, services...
