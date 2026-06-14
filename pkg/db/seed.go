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
		{Name: "Slime Rojo", SpriteKey: "enemy_slime", BaseHP: 45, BaseAttack: 7, BaseXP: 15},
		{Name: "Slime Guerrero", SpriteKey: "enemy_slime", BaseHP: 60, BaseAttack: 10, BaseXP: 25},
		{Name: "Slime Ácido", SpriteKey: "enemy_slime", BaseHP: 50, BaseAttack: 8, BaseXP: 20},
		{Name: "Mini Slime", SpriteKey: "enemy_slime", BaseHP: 20, BaseAttack: 3, BaseXP: 5},
		{Name: "Slime de Cristal", SpriteKey: "enemy_slime", BaseHP: 80, BaseAttack: 12, BaseXP: 40},
	}
	for _, enemy := range enemies {
		DB.Where("name = ?", enemy.Name).FirstOrCreate(&enemy)
	}

	// 2. Sembrar Objetos (Items)
	items := []models.Item{
		{Name: "Espada de Hierro", SpriteKey: "item_sword", Description: "+50% Damage", Type: "Weapon", StatsModifier: `{"damageMultiplier": 1.5}`},
		{Name: "Cerveza de Haste", SpriteKey: "item_beer", Description: "+30% Attack Speed", Type: "Potion", StatsModifier: `{"attackSpeedMultiplier": 1.3}`},
		{Name: "Poción de Salud", SpriteKey: "item_potion", Description: "Restore 50 HP", Type: "Potion", IsConsumable: true, StatsModifier: `{"heal": 50}`},
		{Name: "Poción Pequeña", SpriteKey: "item_potion_small", Description: "Restore 25 HP", Type: "Potion", IsConsumable: true, StatsModifier: `{"heal": 25}`},
		{Name: "Té Enfermizo", SpriteKey: "item_tea", Description: "Velocidad +10%", Type: "Potion", StatsModifier: `{"speedMultiplier": 1.1}`},
		{Name: "Sack of Weight", SpriteKey: "item_sack", Description: "-50% Speed, +50 HP", Type: "Accessory", StatsModifier: `{"speedMultiplier": 0.5, "maxHpBoost": 50}`},
		{Name: "Armadura de Cuero", SpriteKey: "item_armor", Description: "Protección ligera (+10 Vida Máxima).", Type: "Armor", StatsModifier: `{"hp": 10}`},
	}
	for _, item := range items {
		DB.Where("name = ?", item.Name).FirstOrCreate(&item)
	}

	// 3. Sembrar Mapas Prefabricados
	maps := []models.Map{
		{Name: "Nivel 1: Mazmorra de Cristal", Layout: `["####DD####", "#_M______#", "#__T__M__#", "#_r___M__#", "#_____M__#", "#_____P__#", "####DD####"]`, Dictionary: `{"M": "Slime", "D": "door", "T": "loot", "P": "Player", "r": "obstacle"}`},
		{Name: "Nivel 2: Pasadizo Angosto", Layout: `["###DD###", "#M____M#", "#__P___#", "#r_M__r#", "###DD###"]`, Dictionary: `{"M": "Slime Rojo", "D": "door", "P": "Player", "r": "obstacle"}`},
		{Name: "Nivel 3: Almacén Abandonado", Layout: `["##########DD##########", "#M______r____r______M#", "#___T________T_______#", "#_r____M____M____r___#", "#______M_P__M________#", "#_r____M____M____r___#", "#___T________T_______#", "#M______r____r______M#", "##########DD##########"]`, Dictionary: `{"M": "Slime Guerrero", "D": "door", "T": "loot", "P": "Player", "r": "obstacle"}`},
		{Name: "Nivel 4: Ruinas Enredadas", Layout: `["####DD####", "#P_r____M#", "#r_r_M_r_#", "#__M_r_r_#", "#r_r_M___#", "#T___r_M_#", "####DD####"]`, Dictionary: `{"M": "Slime Ácido", "D": "door", "T": "loot", "P": "Player", "r": "obstacle"}`},
		{Name: "Nivel 5: Foso de Práctica", Layout: `["#######", "#M_M_M#", "#_T_T_#", "#M_P_M#", "#_r_r_#", "#M_M_M#", "###D###"]`, Dictionary: `{"M": "Mini Slime", "D": "door", "T": "loot", "P": "Player", "r": "obstacle"}`},
		{Name: "Nivel 6: Galería de los Espejos", Layout: `["######################", "#____r_______r___r___#", "D_P__r_M_r_M_r_M_T_M_D", "#________r_______r___#", "######################"]`, Dictionary: `{"M": "Slime de Cristal", "D": "door", "T": "loot", "P": "Player", "r": "obstacle"}`},
	}
	for _, gameMap := range maps {
		DB.Where("name = ?", gameMap.Name).FirstOrCreate(&gameMap)
	}

	// 4. Sembrar Conocimiento Base (RAG)
	knowledge := []models.KnowledgeChunk{
		{Keywords: `["lore", "mundo"]`, Content: "PrismaCrawler es un laberinto infinito donde los aventureros buscan la redención."},
		{Keywords: `["enemigos", "slime", "goblin"]`, Content: "Los enemigos básicos incluyen Slimes (30 HP) y Goblins (50 HP)."},
	}
	for _, k := range knowledge {
		DB.Where("content = ?", k.Content).FirstOrCreate(&k)
	}

	log.Println("🌱 Datos semilla (Seeder) insertados correctamente.")
}
