package handler

import (
	"encoding/json"
	"net/http"
)

const vsBase = 0xE0100

func fromVariationSelector(r rune) (byte, bool) {
	if r >= vsBase && r < vsBase+256 {
		return byte(r - vsBase), true
	}
	return 0, false
}

func decodeText(text string) string {
	var decoded []byte
	runes := []rune(text)
	for _, r := range runes {
		if b, ok := fromVariationSelector(r); ok {
			decoded = append(decoded, b)
		} else if len(decoded) > 0 {
			break
		}
	}
	return string(decoded)
}

type DecodeRequest struct {
	Text string `json:"text"`
}

type DecodeResponse struct {
	Decoded string `json:"decoded"`
}

func Decode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DecodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	decoded := decodeText(req.Text)
	resp := DecodeResponse{Decoded: decoded}
	json.NewEncoder(w).Encode(resp)
}
