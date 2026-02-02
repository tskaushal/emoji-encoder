package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

const variationSelectorBase = 0xE0100

func toVariationSelector(b byte) rune {
	return rune(variationSelectorBase + int(b))
}

func fromVariationSelector(r rune) (byte, bool) {
	if r >= variationSelectorBase && r < variationSelectorBase+256 {
		return byte(r - variationSelectorBase), true
	}
	return 0, false
}

func Encode(emoji, text string) string {
	var result strings.Builder
	result.WriteString(emoji)
	bytes := []byte(text)
	for _, b := range bytes {
		result.WriteRune(toVariationSelector(b))
	}
	return result.String()
}

func Decode(text string) string {
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

type EncodeRequest struct {
	Emoji string `json:"emoji"`
	Text  string `json:"text"`
}

type EncodeResponse struct {
	Encoded string `json:"encoded"`
}

type DecodeRequest struct {
	Text string `json:"text"`
}

type DecodeResponse struct {
	Decoded string `json:"decoded"`
}

func encodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req EncodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	encoded := Encode(req.Emoji, req.Text)
	resp := EncodeResponse{Encoded: encoded}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req DecodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	decoded := Decode(req.Text)
	resp := DecodeResponse{Decoded: decoded}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/api/encode", encodeHandler)
	http.HandleFunc("/api/decode", decodeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
