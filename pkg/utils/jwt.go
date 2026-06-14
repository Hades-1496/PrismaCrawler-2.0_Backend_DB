package utils

import (
	"errors"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// getSecret recupera el secreto desde la variable de entorno JWT_SECRET.
// Nunca se firma/valida con un secreto débil conocido: la ausencia de JWT_SECRET
// es un error de configuración que se detecta al arrancar (main.go).
func getSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Nunca firmar/validar con un secreto débil conocido. La ausencia de
		// JWT_SECRET es un error de configuración: se valida al arrancar (main.go).
		panic("JWT_SECRET no está configurado")
	}
	return []byte(secret)
}

// GenerateToken crea un JWT para un usuario
func GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // El token expira en 24 horas
	})

	return token.SignedString(getSecret())
}

// ValidateToken verifica el token y extrae los datos (claims)
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getSecret(), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}
