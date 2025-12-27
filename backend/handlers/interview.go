package handlers

import (
	"fmt"
	"net/http"
)

type InterviewRequest struct {
	Question string `json:"question"`
}

func (h *Handler) HandleInterview(w http.ResponseWriter, r *http.Request){
	// Ensure this is a POST request
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}
	var req InterviewRequest
	if err := h.parseJSONBody(r, &req); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Question == "" {
		h.writeError(w, http.StatusBadRequest, "Missing field: question")
		return
	}

	prompt := fmt.Sprintf("You are a DevOps technical interview bot. You are also a DevOps expert with DevOps Interview Knowldege Answer this question clearly & in very detail with examples well elaborated: %s", req.Question)
	explanation, err := h.generateExplanation(r.Context(), prompt)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to generate explanation")
		return
	}

	h.writeJSON(w, http.StatusOK, AIResponse{Prompt: prompt, Explanation: explanation})
	
}