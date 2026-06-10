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
