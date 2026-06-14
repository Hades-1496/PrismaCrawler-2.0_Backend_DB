package main

import (
	"log"
	"os"
	"prismacrawler/internal/handlers"
	"prismacrawler/internal/repository"
	"prismacrawler/internal/routes"
	"prismacrawler/internal/services"
	"prismacrawler/pkg/aiclient"
	"prismacrawler/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil { // nil es null para pointers, maps, etc.
		log.Println("Aviso: No se encontró archivo .env. Se leerán las variables de entorno del sistema.")
	}
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("FATAL: JWT_SECRET no está configurado. Aborta el arranque por seguridad.")
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
	var aiObserver services.GameObserver
	if handlers.AI != nil {
		aiObserver = services.NewAIGameObserver(handlers.AI)
	}
	gameSvc := services.NewGameService(gameRepo, aiObserver)
	gameHandler := handlers.NewGameHandler(gameSvc)
	internalHandler := handlers.NewInternalHandler()

	// Routes
	router := gin.Default()

	// Separación de responsabilidades (SoC): Delegamos TODA la configuración de rutas a su propio paquete.
	routes.Setup(router, gameHandler, internalHandler)

	router.Run(":" + PORT)
}
