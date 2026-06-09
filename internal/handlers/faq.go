package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"prismacrawler/pkg/aiclient"
	"prismacrawler/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AI es el cliente hacia el backend de IA. Se inyecta desde main (ver cmd/api/main.go).
var AI *aiclient.Client

// respondAIError traduce un fallo del cliente de IA a un código HTTP limpio,
// sin filtrar internos (Gemini, Discord, etc.). El detalle solo se loguea.
func respondAIError(c *gin.Context, context string, err error) {
	switch err {
	case aiclient.ErrTimeout:
		utils.SendError(c, http.StatusGatewayTimeout, "El asistente tardó demasiado en responder. Inténtalo de nuevo.")
	case aiclient.ErrBadRequest:
		utils.SendError(c, http.StatusBadRequest, "La petición no pudo procesarse.")
	default: // ErrUnavailable y cualquier otro fallo
		utils.SendError(c, http.StatusBadGateway, "El servicio de IA no está disponible ahora mismo.")
	}
	log.Printf("AI %s error: %v", context, err) // detalle solo en el servidor
}

// Faq proxea la pregunta del chatbot al microservicio de IA.
func Faq(c *gin.Context) {
	var req FaqRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Pregunta no válida (3-500 caracteres).")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	res, err := AI.AskFaq(ctx, req.Question)
	if err != nil {
		respondAIError(c, "faq", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": res.Answer, "sources": res.Sources})
}
