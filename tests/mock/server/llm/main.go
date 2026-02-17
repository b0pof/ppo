package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type ChatRequest struct {
	Message string `json:"message"`
	Stream  bool   `json:"stream"`
}

type ChatResponse struct {
	Text         string `json:"text"`
	GenerationID string `json:"generation_id"`
	ResponseID   string `json:"response_id"`
	FinishReason string `json:"finish_reason"`
	ChatHistory  []struct {
		Message string `json:"message"`
		Role    string `json:"role"`
	} `json:"chat_history"`
	Meta struct {
		APIVersion struct {
			Version string `json:"version"`
		} `json:"api_version"`
		BilledUnits struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"billed_units"`
		Tokens struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"tokens"`
		CachedTokens int `json:"cached_tokens"`
	} `json:"meta"`
}

func main() {
	http.HandleFunc("/v1/chat", chatHandler)

	port := "6666"

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req ChatRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response := ChatResponse{
		Text:         "Тестовый ответ",
		GenerationID: "2c3a5854-95c4-4cc8-8ed9-a50a317705c0",
		ResponseID:   "83f76ff7-03af-4c79-814b-0df8cdee481e",
		FinishReason: "COMPLETE",
		ChatHistory: []struct {
			Message string `json:"message"`
			Role    string `json:"role"`
		}{
			{
				Message: req.Message,
				Role:    "USER",
			},
			{
				Message: "Тестовый ответ",
				Role:    "CHATBOT",
			},
		},
		Meta: struct {
			APIVersion struct {
				Version string `json:"version"`
			} `json:"api_version"`
			BilledUnits struct {
				InputTokens  int `json:"input_tokens"`
				OutputTokens int `json:"output_tokens"`
			} `json:"billed_units"`
			Tokens struct {
				InputTokens  int `json:"input_tokens"`
				OutputTokens int `json:"output_tokens"`
			} `json:"tokens"`
			CachedTokens int `json:"cached_tokens"`
		}{
			APIVersion: struct {
				Version string `json:"version"`
			}{
				Version: "1",
			},
			BilledUnits: struct {
				InputTokens  int `json:"input_tokens"`
				OutputTokens int `json:"output_tokens"`
			}{
				InputTokens:  29,
				OutputTokens: 54,
			},
			Tokens: struct {
				InputTokens  int `json:"input_tokens"`
				OutputTokens int `json:"output_tokens"`
			}{
				InputTokens:  524,
				OutputTokens: 56,
			},
			CachedTokens: 448,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка при отправке ответа: %v", err)
	}
}
