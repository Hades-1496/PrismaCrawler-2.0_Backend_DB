package main

import (
	"log"
	"os"
	"prismacrawler/internal/handlers"
	"prismacrawler/internal/middlewares"
	"prismacrawler/pkg/aiclient"
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
	// DIRECT_URL usa el puerto 5432 (conexión directa), necesario para GORM.
	// DATABASE_URL usa el pooler pgbouncer (6543) que rompe los prepared statements de GORM.
	dbURL := os.Getenv("DIRECT_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	db.ConnectDB(dbURL)

	// Insertamos los datos básicos por defecto (Enemigos, Items, Mapas)
	db.SeedData()

	// Cliente hacia el microservicio de IA. La URL nunca se expone al frontend;
	// el token compartido viaja en X-Internal-Token para que la IA sea "ciega".
	aiURL := os.Getenv("AI_SERVICE_URL")
	if aiURL == "" {
		aiURL = "http://localhost:8001"
	}
	handlers.AI = aiclient.New(aiURL, os.Getenv("AI_INTERNAL_TOKEN"))

	// Routes
	router := gin.Default()

	// Aplicamos CORS de forma global a todas las rutas
	router.Use(middlewares.CORSMiddleware())

	// Leaderboard PÚBLICO (sin JWT): se registra en el router base, fuera del
	// grupo protegido. Es la única fuente del ranking (antes lo servía la IA).
	router.GET("/api/leaderboard", handlers.GetLeaderboard)

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
		apiGroup.GET("/items", handlers.GetItems)

		// Chatbot: el orquestador proxea hacia el backend de IA (la IA es ciega).
		apiGroup.POST("/faq", handlers.Faq)

		// Rutas de Contenido del Juego
		apiGroup.GET("/enemies", handlers.GetEnemies)
		apiGroup.GET("/maps", handlers.GetMaps)
		apiGroup.GET("/maps/:id", handlers.GetMapByID)
	}

	// Grupo de rutas de Administración (Doble protección: Auth + Admin)
	adminGroup := apiGroup.Group("/admin")
	adminGroup.Use(middlewares.AdminMiddleware())
	{
		adminGroup.PUT("/role", handlers.UpdateRole)

		// Acciones de Discord proxeadas a la IA (el front nunca llama a la IA directamente).
		adminGroup.POST("/discord/changelog", handlers.DiscordChangelog)
		adminGroup.POST("/discord/test-webhook", handlers.DiscordTestWebhook)
	}

	router.Run(":" + PORT)
}
