package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"prismacrawler/pkg/utils"
)

// TestGetItems comprueba que el endpoint del catálogo devuelva un 200 OK
func TestGetItems(t *testing.T) {
	// 1. Cargamos el entorno de pruebas seguro (BD paralela y Rutas configuradas)
	router := SetupTestRouter()

	// 2. Generamos un token válido para pasar el AuthMiddleware
	token, _ := utils.GenerateToken(1)

	// 3. Simulamos una petición HTTP GET a la ruta
	req, err := http.NewRequest(http.MethodGet, "/api/items", nil)
	if err != nil {
		t.Fatalf("No se pudo crear la petición: %v", err)
	}

	// 4. Añadimos el token a la cabecera de la petición
	req.Header.Set("Authorization", "Bearer "+token)

	// 5. Creamos un "Grabador de Respuestas" para atrapar lo que devuelve Gin
	w := httptest.NewRecorder()

	// 6. Ejecutamos la petición
	router.ServeHTTP(w, req)

	// 7. Afirmaciones (Asserts)
	// Comprobamos que el código de estado sea 200 OK
	if w.Code != http.StatusOK {
		t.Errorf("Se esperaba código %d pero se obtuvo %d", http.StatusOK, w.Code)
	}
	// Aquí podríamos comprobar también si w.Body.String() contiene "items"
}
