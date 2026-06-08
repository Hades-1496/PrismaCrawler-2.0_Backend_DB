package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"prismacrawler/internal/handlers"
	"prismacrawler/internal/middlewares"
)

// TestUserProfile_Unauthorized comprueba que el middleware bloquea peticiones sin token
func TestUserProfile_Unauthorized(t *testing.T) {
	router := SetupTestRouter()

	apiGroup := router.Group("/api")
	apiGroup.Use(middlewares.AuthMiddleware()) // Aplicamos el middleware
	{
		apiGroup.GET("/profile", handlers.GetProfile)
	}

	req, _ := http.NewRequest(http.MethodGet, "/api/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req) // Hacemos la petición SIN poner la cabecera Authorization

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Seguridad fallida: se esperaba 401 Unauthorized, se obtuvo %d", w.Code)
	}
}