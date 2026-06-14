package handlers

import (
	"encoding/json"
	"hash/fnv"
	"math/rand"
	"net/http"
	"strconv"

	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetEnemies devuelve el catálogo completo de enemigos (bestiario)
func GetEnemies(c *gin.Context) {
	var enemies []models.Enemy
	if err := db.DB.Find(&enemies).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error al obtener el bestiario")
		return
	}
	c.JSON(http.StatusOK, gin.H{"enemies": enemies})
}

// GetMaps devuelve la lista de mapas prefabricados
func GetMaps(c *gin.Context) {
	var maps []models.Map
	if err := db.DB.Find(&maps).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error al obtener la lista de mapas")
		return
	}
	c.JSON(http.StatusOK, gin.H{"maps": maps})
}

// GetMapByID devuelve un mapa por ID. Si no existe en DB genera uno procedural.
func GetMapByID(c *gin.Context) {
	idParam := c.Param("id")

	var gameMap models.Map
	if err := db.DB.First(&gameMap, idParam).Error; err == nil {
		c.JSON(http.StatusOK, gameMap)
		return
	}

	// Tabla vacía o ID no encontrado — generamos el mapa proceduralmente
	id, _ := strconv.Atoi(idParam)
	if id < 1 {
		id = 1
	}
	seed := c.Query("seed")
	layout, dictionary := generateMap(id, seed)

	layoutJSON, _ := json.Marshal(layout)
	dictJSON, _ := json.Marshal(dictionary)

	c.JSON(http.StatusOK, gin.H{
		"ID":         id,
		"Name":       "Sector " + idParam,
		"Level":      id,
		"Layout":     string(layoutJSON),
		"Dictionary": string(dictJSON),
	})
}

// seedSource combina el nivel y la semilla de la run para un RNG determinista por partida.
// Sin seed (cadena vacía) es compatible hacia atrás (solo hash del nivel).
func seedSource(level int, seed string) int64 {
	h := fnv.New64a()
	h.Write([]byte(seed))
	return int64(level*31337) ^ int64(h.Sum64())
}

// generateMap genera un dungeon ASCII determinista basado en el nivel y la semilla de run.
func generateMap(level int, seed string) ([]string, map[string]any) {
	rng := rand.New(rand.NewSource(seedSource(level, seed)))

	cols := 20
	rows := 14

	// Matriz interna de celdas
	grid := make([][]byte, rows)
	for r := range grid {
		grid[r] = make([]byte, cols)
		for c := range grid[r] {
			grid[r][c] = '#'
		}
	}

	// Excavar habitaciones
	type room struct{ x, y, w, h int }
	var rooms []room
	attempts := 0
	for len(rooms) < 6 && attempts < 60 {
		attempts++
		rw := rng.Intn(4) + 3 // ancho 3-6
		rh := rng.Intn(3) + 3 // alto 3-5
		rx := rng.Intn(cols-rw-2) + 1
		ry := rng.Intn(rows-rh-2) + 1

		// Comprobar solapamiento
		overlap := false
		for _, rm := range rooms {
			if rx < rm.x+rm.w+1 && rx+rw > rm.x-1 && ry < rm.y+rm.h+1 && ry+rh > rm.y-1 {
				overlap = true
				break
			}
		}
		if overlap {
			continue
		}
		rooms = append(rooms, room{rx, ry, rw, rh})
		for dy := 0; dy < rh; dy++ {
			for dx := 0; dx < rw; dx++ {
				grid[ry+dy][rx+dx] = '_'
			}
		}
	}

	// Conectar habitaciones con pasillos
	for i := 1; i < len(rooms); i++ {
		a := rooms[i-1]
		b := rooms[i]
		ax, ay := a.x+a.w/2, a.y+a.h/2
		bx, by := b.x+b.w/2, b.y+b.h/2
		for x := min(ax, bx); x <= max(ax, bx); x++ {
			grid[ay][x] = '_'
		}
		for y := min(ay, by); y <= max(ay, by); y++ {
			grid[y][bx] = '_'
		}
	}

	// Colocar jugador en la primera habitación
	if len(rooms) > 0 {
		r := rooms[0]
		grid[r.y+r.h/2][r.x+1] = 'P'
	}

	// Colocar salida 'E' SIEMPRE (incluso con una sola sala), lejos de 'P'.
	if len(rooms) > 0 {
		last := rooms[len(rooms)-1]
		ex, ey := last.x+last.w-2, last.y+last.h/2
		if grid[ey][ex] != '_' { // chocaría con 'P' o un muro: buscar otra celda de suelo
		search:
			for r := rows - 1; r >= 0; r-- {
				for c := cols - 1; c >= 0; c-- {
					if grid[r][c] == '_' {
						ex, ey = c, r
						break search
					}
				}
			}
		}
		grid[ey][ex] = 'E'
	}

	// Colocar monstruos sobre celdas de suelo libres (ACOTADO: nunca bucle infinito).
	monsterCount := 2 + level/2
	if monsterCount > 8 {
		monsterCount = 8
	}
	type cell struct{ r, c int }
	var freeCells []cell
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == '_' {
				freeCells = append(freeCells, cell{r, c})
			}
		}
	}
	// Barajado determinista con el RNG de la run; placeable = min(deseado, disponible).
	rng.Shuffle(len(freeCells), func(i, j int) {
		freeCells[i], freeCells[j] = freeCells[j], freeCells[i]
	})
	if monsterCount > len(freeCells) {
		monsterCount = len(freeCells)
	}
	for i := 0; i < monsterCount; i++ {
		grid[freeCells[i].r][freeCells[i].c] = 'M'
	}

	// Serializar a slice de strings
	result := make([]string, rows)
	for r := range grid {
		result[r] = string(grid[r])
	}

	hp := 20 + level*10
	dmg := 5 + level*3
	dictionary := map[string]any{
		"M": map[string]any{"hp": hp, "damage": dmg},
	}

	return result, dictionary
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
