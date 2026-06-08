package db

import (
	"log"
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

	log.Println("!Conectado a la base de datos exitosamente!")

	// GORM creará/actualizará las tablas basándose en las structs
	// Migramos los 5 modelos principales del MVP
	DB.AutoMigrate(&models.User{}, &models.Character{}, &models.Item{}, &models.GameRun{}, &models.RunInventory{})
}
