package db

import (
	"log"
	"time"

	"prismacrawler/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(url string) {
	var err error
	DB, err = gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatal("Fallo al contactar a la base de datos: ", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Error obteniendo instancia sql.DB: ", err)
	}

	// Evita conexiones muertas: rota antes de que Supabase las cierre por inactividad
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	log.Println("!Conectado a la base de datos exitosamente!")

	// Habilitamos la extensión pgvector para el microservicio de IA
	DB.Exec("CREATE EXTENSION IF NOT EXISTS vector;")

	// GORM creará/actualizará las tablas basándose en las structs
	// Migramos los modelos del juego y los modelos de soporte para la IA
	DB.AutoMigrate(&models.User{}, &models.Character{}, &models.Item{}, &models.GameRun{}, &models.RunInventory{}, &models.Map{}, &models.Enemy{}, &models.KnowledgeChunk{}, &models.Wallet{}, &models.Garden{}, &models.WeeklyReward{}, &models.ProcessedStripeEvent{})
}
