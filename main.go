package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	//go:embed static/*
	staticFS embed.FS
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
	for _, b := range []byte(text) {
		result.WriteRune(toVariationSelector(b))
	}
	return result.String()
}

func Decode(text string) string {
	var decoded []byte
	for _, r := range []rune(text) {
		if b, ok := fromVariationSelector(r); ok {
			decoded = append(decoded, b)
		}
	}
	return string(decoded)
}

func encodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Emoji string `json:"emoji"`
		Text  string `json:"text"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	json.NewEncoder(w).Encode(map[string]string{
		"encoded": Encode(req.Emoji, req.Text),
	})
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Text string `json:"text"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	json.NewEncoder(w).Encode(map[string]string{
		"decoded": Decode(req.Text),
	})
}

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.FS(staticFS)),
		),
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, _ := staticFS.ReadFile("static/index.html")
		w.Write(data)
	})

	http.HandleFunc("/api/encode", encodeHandler)
	http.HandleFunc("/api/decode", decodeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
