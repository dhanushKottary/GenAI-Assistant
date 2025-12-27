package handlers

import (
	"fmt"
	"net/http"
)

type CodeRequest struct {
	Code string `json:"code"`
}



func (h *Handler) CodeExplainer(w http.ResponseWriter, r *http.Request) {
	// Ensure this is a POST request
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	// Parse JSON body
	var req CodeRequest
	if err := h.parseJSONBody(r, &req); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Code == "" {
		h.writeError(w, http.StatusBadRequest, "Missing field: code")
		return
	}
	prompt := fmt.Sprintf("You are a DevOps Expert & I want You to Explain this code in detail with examples & real scenarios:\n\n%s", req.Code)
	explanation, err := h.generateExplanation(r.Context(), prompt)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to generate explanation")
		return
	}

    h.writeJSON(w, http.StatusOK, AIResponse{Prompt: prompt, Explanation: explanation})
}
