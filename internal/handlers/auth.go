package handlers

import (
	"context"
	"log"
	"net/http"
	"time"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register maneja el registro de nuevos usuarios
func Register(c *gin.Context) {
	var req RegisterRequest

	// 1. Recibir y validar el JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// 2. Encriptar la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error interno al procesar la contraseña")
		return
	}

	// 3. Crear el modelo del nuevo usuario
	var count int64
	db.DB.Model(&models.User{}).Count(&count)
	// Solo el PRIMER usuario registrado es ADMIN (bootstrap). El resto se
	// promueve de forma segura por un admin existente vía PUT /api/admin/role.
	role := "USER"
	if count == 0 {
		role = "ADMIN"
	}

	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		Nickname:     req.Nickname,
	}

	// 4. Guardar en la base de datos
	result := db.DB.Create(&user)
	if result.Error != nil {
		utils.SendError(c, http.StatusConflict, "Este email o nickname ya está en uso")
		return
	}

	// 5. Generamos el token JWT (Auto-login equivalente al de JS)
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error interno al generar el token")
		return
	}

	// Disparar evento asíncrono para n8n/IA (Goroutine)
	go func(userID uint, email, nickname string) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		payload := gin.H{
			"event": "user_registered",
			"payload": gin.H{
				"user_id":  userID,
				"email":    email,
				"username": nickname,
			},
		}
		if AI != nil {
			if _, err := AI.Proxy(ctx, "/api/n8n/relay", payload); err != nil {
				log.Printf("Error enviando trigger user_registered a IA: %v", err)
			}
		}
	}(user.ID, user.Email, user.Nickname)

	// 5. Responder con éxito (HTTP 201 - Created)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario registrado exitosamente",
		"user_id": user.ID,
		"token":   token,
	})
}

// Login maneja la autenticación y devuelve un JWT
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	var user models.User
	// 1. Buscamos al usuario en base al correo (Equivalente a prisma.user.findUnique)
	result := db.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		utils.SendError(c, http.StatusUnauthorized, "Credenciales inválidas")
		return
	}

	// 2. Comparamos la contraseña encriptada (Equivalente a bcrypt.compare)
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, "Credenciales inválidas")
		return
	}

	// 3. Generamos el token JWT
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error interno al generar el token")
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user_id": user.ID, "email": user.Email, "role": user.Role, "nickname": user.Nickname})
}
