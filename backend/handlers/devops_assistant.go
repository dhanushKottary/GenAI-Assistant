package handlers

import (
	"fmt"
	"net/http"
)

type ReviewRequest struct {
	Content string `json:"content"`
}

func (h *Handler) ReviewDevOps(w http.ResponseWriter, r *http.Request){
	// Ensure this is a POST request
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}
	var req ReviewRequest
	if err := h.parseJSONBody(r, &req); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Content == "" {
		h.writeError(w, http.StatusBadRequest, "Missing field: content")
		return
	}

	prompt := fmt.Sprintf("You are a DevOps expert. I am stuck with this error, I need your expertise to troubleshoot & resolve this error. Explain in very detail:\n\n%s", req.Content)
	explanation, err := h.generateExplanation(r.Context(), prompt)
	
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to generate explanation")
		return
	}
    h.writeJSON(w, http.StatusOK, AIResponse{Prompt: prompt, Explanation: explanation})
	
}