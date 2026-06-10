package tests

import (
	"log"
	"os"
	"prismacrawler/internal/handlers"
	"prismacrawler/internal/models"
	"prismacrawler/internal/repository"
	"prismacrawler/internal/routes"
	"prismacrawler/internal/services"
	"prismacrawler/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// SetupTestRouter configura Gin y una base de datos aislada para los tests
func SetupTestRouter() *gin.Engine {
	// 1. Cargamos el .env desde el directorio padre (la raíz)
	err := godotenv.Load("../.env.test")
	if err != nil {
		log.Println("Aviso: No se encontró archivo .env.test. Asegúrate de tener las variables configuradas.")
	}

	// 2. Usamos la variable ESPECÍFICA para tests para no borrar la DB real
	testDBUrl := os.Getenv("DIRECT_URL")
	if testDBUrl == "" {
		testDBUrl = os.Getenv("DATABASE_URL")
	}
	if testDBUrl == "" {
		log.Fatal("¡ALTO! Necesitas definir DIRECT_URL o DATABASE_URL en tu .env.test para correr los tests de forma segura.")
	}

	// 3. Conectamos GORM a la base de datos de test y reconstruimos las tablas
	db.ConnectDB(testDBUrl)
	db.DB.Exec("DROP TABLE IF EXISTS run_inventories, game_runs, characters, items, users, maps, enemies, knowledge_chunks CASCADE;")
	db.DB.AutoMigrate(&models.User{}, &models.Character{}, &models.Item{}, &models.GameRun{}, &models.RunInventory{}, &models.Map{}, &models.Enemy{}, &models.KnowledgeChunk{})

	// 4. Configuramos el router de prueba
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// 5. Inyectamos dependencias (dejando la IA en nil para no disparar webhooks)
	gameRepo := repository.NewGameRepository(db.DB)
	gameSvc := services.NewGameService(gameRepo, nil)
	gameHandler := handlers.NewGameHandler(gameSvc)
	routes.Setup(router, gameHandler)

	return router
}
