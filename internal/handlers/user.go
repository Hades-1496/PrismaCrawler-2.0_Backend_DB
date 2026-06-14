package handlers

import (
	"net/http"
	"os"
	"strings"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetProfile devuelve los datos del usuario, su billetera y su Top 5 de partidas
func GetProfile(c *gin.Context) {
	// 1. Obtener el ID del usuario desde el JWT
	userID := utils.GetUserID(c)

	// 2. Buscar al usuario
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		utils.SendError(c, http.StatusNotFound, "Usuario no encontrado")
		return
	}

	// 3. Buscar el Top 5 de mejores partidas DE ESTE USUARIO
	// Para ello tenemos que hacer un JOIN con la tabla Characters
	var runs []models.GameRun
	db.DB.Preload("Character").
		Joins("JOIN characters ON characters.id = game_runs.character_id").
		Where("characters.user_id = ? AND game_runs.score > 0", userID).
		Order("game_runs.score desc, game_runs.current_floor desc").
		Limit(5).
		Find(&runs)

	// 4. Mapear las partidas para limpiar la respuesta
	var topRuns []gin.H
	for _, run := range runs {
		topRuns = append(topRuns, gin.H{
			"character": run.Character.Name,
			"class":     run.Character.Class,
			"score":     run.Score,
			"floor":     run.CurrentFloor,
		})
	}

	// 5. Buscar la billetera del usuario
	var wallet models.Wallet
	db.DB.Where("user_id = ?", userID).FirstOrCreate(&wallet, models.Wallet{UserID: userID})

	// 6. Devolver el JSON
	c.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"email":       user.Email,
		"role":        user.Role,
		"nickname":    user.Nickname,
		"real_name":   user.RealName,
		"avatar":      user.Avatar,
		"player_icon": user.PlayerIcon,
		"created_at":  user.CreatedAt,
		"last_login":  user.LastLogin,
		"wallet": gin.H{
			"coins": wallet.Coins,
			"gems":  wallet.Gems,
		},
		"top_runs": topRuns,
	})
}

// UpdateProfile actualiza la información personal del usuario
func UpdateProfile(c *gin.Context) {
	userID := utils.GetUserID(c)
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Datos de perfil inválidos: "+err.Error())
		return
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		utils.SendError(c, http.StatusNotFound, "Usuario no encontrado")
		return
	}

	user.Nickname = req.Nickname
	user.RealName = req.RealName
	user.Avatar = req.Avatar
	user.PlayerIcon = req.PlayerIcon

	// En local/desarrollo permitimos auto-promoción a ADMIN para facilitar pruebas del backoffice
	if req.Role != "" {
		host := c.Request.Host
		env := os.Getenv("ENVIRONMENT")
		if strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1") || env == "development" || env == "" {
			if req.Role == "ADMIN" || req.Role == "USER" {
				user.Role = req.Role
			}
		}
	}

	if err := db.DB.Save(&user).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error al guardar el perfil")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Perfil actualizado exitosamente",
		"nickname":    user.Nickname,
		"real_name":   user.RealName,
		"avatar":      user.Avatar,
		"player_icon": user.PlayerIcon,
		"role":        user.Role,
	})
}

// UpdateRole cambia el rol de un usuario (debería estar protegida por AdminMiddleware)
func UpdateRole(c *gin.Context) {
	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	var user models.User
	if err := db.DB.First(&user, req.UserID).Error; err != nil {
		utils.SendError(c, http.StatusNotFound, "Usuario a modificar no encontrado")
		return
	}

	// Actualizamos el rol y guardamos
	user.Role = req.Role
	if err := db.DB.Save(&user).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error al actualizar el rol del usuario")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rol actualizado exitosamente", "user_id": user.ID, "new_role": user.Role})
}

// ListUsers devuelve la lista de todos los usuarios registrados
func ListUsers(c *gin.Context) {
	var users []models.User
	if err := db.DB.Order("created_at desc").Find(&users).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Error al listar usuarios")
		return
	}

	var response []gin.H
	for _, u := range users {
		var wallet models.Wallet
		db.DB.Where("user_id = ?", u.ID).FirstOrCreate(&wallet, models.Wallet{UserID: u.ID})

		response = append(response, gin.H{
			"id":         u.ID,
			"email":      u.Email,
			"role":       u.Role,
			"nickname":   u.Nickname,
			"real_name":  u.RealName,
			"created_at": u.CreatedAt,
			"wallet": gin.H{
				"coins": wallet.Coins,
				"gems":  wallet.Gems,
			},
		})
	}

	c.JSON(http.StatusOK, response)
}
