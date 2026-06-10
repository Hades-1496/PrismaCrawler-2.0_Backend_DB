package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestEndToEndFlow simula el ciclo completo de un jugador nuevo
func TestEndToEndFlow(t *testing.T) {
	router := SetupTestRouter()

	// PASO 1 y 2: Registro y Login
	userPayload := []byte(`{"email": "e2e@test.com", "password": "password123"}`)

	reqReg, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(userPayload))
	router.ServeHTTP(httptest.NewRecorder(), reqReg) // Ejecutamos registro y lo ignoramos

	reqLog, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(userPayload))
	wLog := httptest.NewRecorder()
	router.ServeHTTP(wLog, reqLog)

	// Extraemos el Token JWT de la respuesta del login
	var loginResp map[string]interface{}
	json.Unmarshal(wLog.Body.Bytes(), &loginResp)
	token := loginResp["token"].(string)

	// PASO 3: Crear el Personaje
	charPayload := []byte(`{"name": "E2E Hero", "class": "Guerrero"}`)
	reqChar, _ := http.NewRequest(http.MethodPost, "/api/characters", bytes.NewBuffer(charPayload))
	reqChar.Header.Set("Authorization", "Bearer "+token) // Añadimos el JWT al Header
	wChar := httptest.NewRecorder()
	router.ServeHTTP(wChar, reqChar)

	if wChar.Code != http.StatusCreated {
		t.Fatalf("Fallo al crear personaje E2E: %d", wChar.Code)
	}

	// Extraemos el ID del personaje recién creado
	var charResp map[string]interface{}
	json.Unmarshal(wChar.Body.Bytes(), &charResp)
	charData := charResp["character"].(map[string]interface{})
	charID := int(charData["ID"].(float64))

	// PASO 4: Iniciar una Partida
	runPayload, _ := json.Marshal(map[string]int{"character_id": charID}) // Construimos el JSON dinámicamente
	reqRun, _ := http.NewRequest(http.MethodPost, "/api/runs/start", bytes.NewBuffer(runPayload))
	reqRun.Header.Set("Authorization", "Bearer "+token)
	wRun := httptest.NewRecorder()
	router.ServeHTTP(wRun, reqRun)

	if wRun.Code != http.StatusCreated {
		t.Errorf("Fallo al iniciar partida E2E: %d", wRun.Code)
	}
}
