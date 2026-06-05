package db


import (
	"log"
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
}