package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

const variationSelectorBase = 0xE0100

func toVariationSelector(b byte) rune {
	return rune(variationSelectorBase + int(b))
}

func encodeText(emoji, text string) string {
	var result strings.Builder
	result.WriteString(emoji)
	bytes := []byte(text)
	for _, b := range bytes {
		result.WriteRune(toVariationSelector(b))
	}
	return result.String()
}

type EncodeRequest struct {
	Emoji string `json:"emoji"`
	Text  string `json:"text"`
}

type EncodeResponse struct {
	Encoded string `json:"encoded"`
}

func Encode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req EncodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	encoded := encodeText(req.Emoji, req.Text)
	resp := EncodeResponse{Encoded: encoded}
	json.NewEncoder(w).Encode(resp)
}
