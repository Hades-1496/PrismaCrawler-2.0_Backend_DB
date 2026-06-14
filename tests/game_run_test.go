package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"prismacrawler/internal/handlers"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveRun_WithItems(t *testing.T) {
	router := SetupTestRouter()

	user := models.User{Email: "saverun@test.com", Nickname: "test"}
	db.DB.Create(&user)

	char := models.Character{UserID: user.ID, Name: "Hero", Class: "Warrior", BaseHP: 100, IsAlive: true}
	db.DB.Create(&char)

	item := models.Item{Name: "Sword", SpriteKey: "item_sword", Type: "Weapon", StatsModifier: "{}"}
	db.DB.Create(&item)

	run := models.GameRun{CharacterID: char.ID, Status: "In_Progress"}
	db.DB.Create(&run)

	reqBody := handlers.SaveRunRequest{
		RunID:          run.ID,
		CurrentFloor:   2,
		Score:          100,
		CurrentHP:      80,
		Kills:          5,
		DamageDealt:    200,
		DamageTaken:    20,
		ItemsCollected: []string{"item_sword"},
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/api/runs/save", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	// Utilizar jwt.go de tu proyecto para autorizar
	token, _ := utils.GenerateToken(user.ID)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var runInv models.RunInventory
	err := db.DB.Where("run_id = ? AND item_id = ?", run.ID, item.ID).First(&runInv).Error
	assert.NoError(t, err)
	assert.Equal(t, 1, runInv.Quantity)
}