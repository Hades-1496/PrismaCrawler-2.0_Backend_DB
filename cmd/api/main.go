package main

import (
	"log"
	"os"
	"prismacrawler/internal/handlers"
	"prismacrawler/internal/repository"
	"prismacrawler/internal/services"
	"prismacrawler/internal/routes"
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

	// Inyección de Dependencias (Wiring)
	// Creamos las instancias de cada capa, inyectando sus dependencias.
	gameRepo := repository.NewGameRepository(db.DB)
	gameSvc := services.NewGameService(gameRepo, handlers.AI)
	gameHandler := handlers.NewGameHandler(gameSvc)

	// Routes
	router := gin.Default()

	// Separación de responsabilidades (SoC): Delegamos la configuración de rutas
	// y le pasamos los handlers que necesitan dependencias.
	routes.Setup(router, gameHandler)

	router.Run(":" + PORT)
}
