package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAuthEndpoints verifica el registro y login de forma aislada
func TestAuthEndpoints(t *testing.T) {
	router := SetupTestRouter()

	// 1. Probamos el Registro
	registerPayload := []byte(`{"email": "hero@test.com", "password": "password123"}`)
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
