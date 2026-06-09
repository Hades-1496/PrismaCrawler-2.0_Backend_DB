package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"prismacrawler/internal/handlers"
)

// TestGetItems comprueba que el endpoint del catálogo devuelva un 200 OK
func TestGetItems(t *testing.T) {
	// 1. Cargamos el entorno de pruebas seguro (Base de datos paralela + Router vacío)
	router := SetupTestRouter()

	// 2. Registramos solo la ruta que queremos testear en este archivo
	router.GET("/api/items", handlers.GetItems)

	// 3. Simulamos una petición HTTP GET a la ruta
	req, err := http.NewRequest(http.MethodGet, "/api/items", nil)
	if err != nil {
		t.Fatalf("No se pudo crear la petición: %v", err)
	}

	// 4. Creamos un "Grabador de Respuestas" para atrapar lo que devuelve Gin
	w := httptest.NewRecorder()

	// 5. Ejecutamos la petición
	router.ServeHTTP(w, req)

	// 6. Afirmaciones (Asserts)
	// Comprobamos que el código de estado sea 200 OK
	if w.Code != http.StatusOK {
		t.Errorf("Se esperaba código %d pero se obtuvo %d", http.StatusOK, w.Code)
	}
	// Aquí podríamos comprobar también si w.Body.String() contiene "items"
}
