package aiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Errores tipados para que los handlers los mapeen a códigos HTTP.
var (
	ErrTimeout     = errors.New("ai: timeout")     // → 504
	ErrUnavailable = errors.New("ai: unavailable") // → 502
	ErrBadRequest  = errors.New("ai: bad request") // → 400/422
)

// Client es el puente interno del orquestador Go hacia el backend de IA.
type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func New(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		http:    &http.Client{Timeout: 20 * time.Second}, // colchón sobre el ctx del handler
	}
}

// post realiza un POST interno a la IA inyectando el token compartido y
// normaliza los fallos a los errores tipados. Devuelve el cuerpo crudo en 2xx.
func (c *Client) post(ctx context.Context, path string, payload any) (json.RawMessage, error) {
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-Token", c.token) // microservicio ciego: solo el orquestador entra

	resp, err := c.http.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, ErrTimeout
		}
		return nil, fmt.Errorf("%w: %v", ErrUnavailable, err) // conexión rechazada, DNS, etc.
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("%w: respuesta ilegible", ErrUnavailable)
		}
		return raw, nil
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		return nil, ErrBadRequest
	default: // 5xx, incluido el 503 de Gemini o el 502 de Discord caídos
		return nil, ErrUnavailable
	}
}

// Proxy reenvía un POST interno arbitrario (p. ej. acciones de Discord) y
// devuelve el cuerpo de respuesta de la IA tal cual para retransmitirlo.
func (c *Client) Proxy(ctx context.Context, path string, payload any) (json.RawMessage, error) {
	return c.post(ctx, path, payload)
}
