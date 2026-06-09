package services

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type aiEvent struct {
	Event   string                 `json:"event"`
	Payload map[string]interface{} `json:"payload"`
}

// NotifyAI envía un evento al backend de IA de forma asíncrona (fire-and-forget).
// Llámalo siempre en una goroutine: go services.NotifyAI(...)
func NotifyAI(event string, payload map[string]interface{}) {
	baseURL := os.Getenv("AI_BACKEND_URL")
	token := os.Getenv("AI_INTERNAL_TOKEN")
	if baseURL == "" {
		return
	}

	data, err := json.Marshal(aiEvent{Event: event, Payload: payload})
	if err != nil {
		log.Printf("[AI Notifier] marshal error para evento %s: %v", event, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/api/n8n/relay", bytes.NewReader(data))
	if err != nil {
		log.Printf("[AI Notifier] error creando request para evento %s: %v", event, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-Token", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[AI Notifier] error enviando evento %s: %v", event, err)
		return
	}
	defer resp.Body.Close()
	log.Printf("[AI Notifier] evento '%s' enviado — HTTP %d", event, resp.StatusCode)
}
