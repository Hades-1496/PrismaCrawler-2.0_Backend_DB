package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// DiscordChangelog proxea la publicación de un changelog hacia el backend de IA
// (módulo Discord). Es una operación de administración: el orquestador ya validó
// auth + rol admin antes de llegar aquí.
func DiscordChangelog(c *gin.Context) {
	var req ChangelogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Changelog no válido (contenido requerido, máx. 4000 caracteres)."})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	raw, err := AI.Proxy(ctx, "/api/discord/changelog", gin.H{"title": req.Title, "content": req.Content})
	if err != nil {
		respondAIError(c, "discord/changelog", err)
		return
	}
	c.Data(http.StatusOK, "application/json", raw)
}

// DiscordTestWebhook proxea el envío de un mensaje de prueba al canal de Discord.
func DiscordTestWebhook(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	raw, err := AI.Proxy(ctx, "/api/discord/test-webhook", json.RawMessage("{}"))
	if err != nil {
		respondAIError(c, "discord/test-webhook", err)
		return
	}
	c.Data(http.StatusOK, "application/json", raw)
}
