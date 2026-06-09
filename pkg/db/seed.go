package db

import (
	"log"
	"prismacrawler/internal/models"
)

// SeedData inserta los valores básicos en la base de datos si no existen
func SeedData() {
	// 1. Sembrar Enemigos
	enemies := []models.Enemy{
		{Name: "Slime", SpriteKey: "enemy_slime", BaseHP: 30, BaseAttack: 5, BaseXP: 10},
		{Name: "Goblin", SpriteKey: "enemy_goblin", BaseHP: 50, BaseAttack: 8, BaseXP: 20},
		{Name: "Orco Oscuro", SpriteKey: "enemy_orc", BaseHP: 100, BaseAttack: 15, BaseXP: 50},
	}
	for _, enemy := range enemies {
		DB.Where("name = ?", enemy.Name).FirstOrCreate(&enemy)
	}

	// 2. Sembrar Objetos (Items)
	items := []models.Item{
		{Name: "Espada de Hierro", SpriteKey: "item_sword", Description: "Una espada básica (+5 Daño).", Type: "Weapon", StatsModifier: `{"damage": 5}`},
		{Name: "Poción de Salud", SpriteKey: "item_potion", Description: "Restaura 50 HP.", Type: "Potion", IsConsumable: true, StatsModifier: `{"heal": 50}`},
		{Name: "Armadura de Cuero", SpriteKey: "item_armor", Description: "Protección ligera (+10 Vida Máxima).", Type: "Armor", StatsModifier: `{"hp": 10}`},
	}
	for _, item := range items {
		DB.Where("name = ?", item.Name).FirstOrCreate(&item)
	}

	// 3. Sembrar Mapas Prefabricados
	maps := []models.Map{
		{Name: "Sala de Entrenamiento", Layout: `["#####", "#P__#", "#_E_#", "#####"]`, Dictionary: `{"P": "player", "E": "Slime"}`},
		{Name: "Emboscada Goblin", Layout: `["#######", "#P_E_E#", "#######"]`, Dictionary: `{"P": "player", "E": "Goblin"}`},
	}
	for _, gameMap := range maps {
		DB.Where("name = ?", gameMap.Name).FirstOrCreate(&gameMap)
	}

	log.Println("🌱 Datos semilla (Seeder) insertados correctamente.")
}
