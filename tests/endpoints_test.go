package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"prismacrawler/internal/models"
	"testing"
)

// TestAuthEndpoints verifica el registro y login de forma aislada
func TestAuthEndpoints(t *testing.T) {
	router := SetupTestRouter()

	// 1. Probamos el Registro
	registerPayload := []byte(`{"email": "hero@test.com", "password": "password123", "nickname": "heroTest"}`)
	reqReg, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(registerPayload))
	reqReg.Header.Set("Content-Type", "application/json") // ¡Importante decirle a Gin que es un JSON!

	wReg := httptest.NewRecorder()
	router.ServeHTTP(wReg, reqReg)

	if wReg.Code != http.StatusCreated {
		t.Errorf("Error en registro: se esperaba 201, se obtuvo %d", wReg.Code)
	}

	// 2. Probamos el Login con las mismas credenciales
	reqLog, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(registerPayload))
	reqLog.Header.Set("Content-Type", "application/json")

	wLog := httptest.NewRecorder()
	router.ServeHTTP(wLog, reqLog)

	if wLog.Code != http.StatusOK {
		t.Errorf("Error en login: se esperaba 200, se obtuvo %d", wLog.Code)
	}
}

// TestErrorHandling verifica el manejo de errores HTTP (400, 401, 404)
func TestErrorHandling(t *testing.T) {
	router := SetupTestRouter()

	// 1. Probamos 400 Bad Request (JSON inválido o campos faltantes)
	badPayload := []byte(`{"email": "not-an-email"}`) // Falta password, email inválido
	reqBad, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(badPayload))
	reqBad.Header.Set("Content-Type", "application/json")

	wBad := httptest.NewRecorder()
	router.ServeHTTP(wBad, reqBad)

	if wBad.Code != http.StatusBadRequest {
		t.Errorf("Error en manejo de 400: se esperaba 400, se obtuvo %d", wBad.Code)
	}

	// 2. Probamos 401 Unauthorized (Login con credenciales incorrectas)
	wrongLogin := []byte(`{"email": "hero@test.com", "password": "wrongpassword"}`)
	reqLog, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(wrongLogin))
	reqLog.Header.Set("Content-Type", "application/json")

	wLog := httptest.NewRecorder()
	router.ServeHTTP(wLog, reqLog)

	if wLog.Code != http.StatusUnauthorized {
		t.Errorf("Error en manejo de 401: se esperaba 401, se obtuvo %d", wLog.Code)
	}

	// 3. Probamos 404 Not Found (Ruta inexistente)
	req404, _ := http.NewRequest(http.MethodGet, "/api/ruta-que-no-existe", nil)
	w404 := httptest.NewRecorder()
	router.ServeHTTP(w404, req404)

	if w404.Code != http.StatusNotFound {
		t.Errorf("Error en manejo de 404: se esperaba 404, se obtuvo %d", w404.Code)
	}
}

// TestPingEndpoint verifica que la ruta de health check responde correctamente
func TestPingEndpoint(t *testing.T) {
	router := SetupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Error en Ping: se esperaba 200, se obtuvo %d", w.Code)
	}
}

// TestIntegralFlow simula el ciclo de vida completo de un usuario en el juego (E2E)
func TestIntegralFlow(t *testing.T) {
	router := SetupTestRouter()

	// 1. Registro
	registerPayload := []byte(`{"email": "integral@test.com", "password": "password123", "nickname": "integralTest"}`)
	reqReg, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(registerPayload))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	router.ServeHTTP(wReg, reqReg)

	if wReg.Code != http.StatusCreated {
		t.Fatalf("Fallo en registro: %d", wReg.Code)
	}

	// 2. Login
	reqLog, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(registerPayload))
	reqLog.Header.Set("Content-Type", "application/json")
	wLog := httptest.NewRecorder()
	router.ServeHTTP(wLog, reqLog)

	if wLog.Code != http.StatusOK {
		t.Fatalf("Fallo en login: %d", wLog.Code)
	}

	var loginResponse struct {
		Token string `json:"token"`
	}
	json.Unmarshal(wLog.Body.Bytes(), &loginResponse)
	token := loginResponse.Token

	// 3. Crear Personaje
	charPayload := []byte(`{"name": "HeroeTest", "class": "Guerrero"}`)
	reqChar, _ := http.NewRequest(http.MethodPost, "/api/characters", bytes.NewBuffer(charPayload))
	reqChar.Header.Set("Content-Type", "application/json")
	reqChar.Header.Set("Authorization", "Bearer "+token)
	wChar := httptest.NewRecorder()
	router.ServeHTTP(wChar, reqChar)

	if wChar.Code != http.StatusCreated {
		t.Fatalf("Fallo en creación de personaje: %d", wChar.Code)
	}

	var charResponse struct {
		Character models.Character `json:"character"`
	}
	json.Unmarshal(wChar.Body.Bytes(), &charResponse)
	charID := charResponse.Character.ID

	// 4. Iniciar Partida
	runPayload, _ := json.Marshal(map[string]uint{"character_id": charID})
	reqRun, _ := http.NewRequest(http.MethodPost, "/api/runs/start", bytes.NewBuffer(runPayload))
	reqRun.Header.Set("Content-Type", "application/json")
	reqRun.Header.Set("Authorization", "Bearer "+token)
	wRun := httptest.NewRecorder()
	router.ServeHTTP(wRun, reqRun)

	if wRun.Code != http.StatusCreated {
		t.Fatalf("Fallo al iniciar partida: %d", wRun.Code)
	}
}
