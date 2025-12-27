package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"google.golang.org/genai"
)

type Handler struct {
	Client *genai.Client
	Model string
	Schema genai.Schema
	Timeout time.Duration
	MaxBodySize int64
}

// NewHandler initializes a reusable handler with sane defaults.
func NewHandler(client *genai.Client) *Handler {
	return &Handler{
		Client: client,
		Model:  "gemini-2.5-flash",
		Schema: genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"explanation": {
					Type:        genai.TypeString,
					Description: "A clear and concise explanation.",
				},
			},
			Required: []string{"explanation"},
		},
		Timeout:     30 * time.Second,
		MaxBodySize: 1 << 20, // 1 MB
	}
}

type AIResponse struct {
	Prompt      string `json:"prompt"`
	Explanation string `json:"explanation"`
}


// --- Utility functions ---

func (h *Handler) parseJSONBody(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return errors.New("empty request body")
	}
	defer r.Body.Close()

	limited := io.LimitReader(r.Body, h.MaxBodySize)
	dec := json.NewDecoder(limited)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, msg string) {
	http.Error(w, msg, status)
}

func (h *Handler) generateExplanation(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, h.Timeout)
	defer cancel()

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   &h.Schema,
	}

	result, err := h.Client.Models.GenerateContent(ctx, h.Model, genai.Text(prompt), config)
	if err != nil {
		return "", fmt.Errorf("Gemini API error: %w", err)
	}

	raw := result.Text()
	var out struct {
		Explanation string `json:"explanation"`
	}
    err = json.Unmarshal([]byte(raw), &out);
	if err != nil {
		log.Printf("unmarshal error: %v | raw: %s", err, raw)
		return "", errors.New("invalid model response")
	}
	return out.Explanation, nil
}