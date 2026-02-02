package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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

const htmlTemplate = `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Hide a message in an emoji</title><style>*{margin:0;padding:0;box-sizing:border-box}body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background:linear-gradient(135deg,#e94ca6 0%,#d946a3 100%);min-height:100vh;display:flex;align-items:center;justify-content:center;padding:20px}.container{background:#fef7f9;border-radius:20px;box-shadow:0 20px 60px rgba(0,0,0,0.3);max-width:500px;width:100%;padding:40px}h1{color:#333;margin-bottom:10px;font-size:1.8em}.description{color:#666;margin-bottom:30px;line-height:1.6}.toggle-container{display:flex;justify-content:center;align-items:center;gap:15px;margin-bottom:30px}.toggle-label{font-weight:600;color:#555}.toggle-switch{position:relative;width:60px;height:30px;background:#ccc;border-radius:15px;cursor:pointer;transition:background 0.3s}.toggle-switch.active{background:#e94ca6}.toggle-slider{position:absolute;top:3px;left:3px;width:24px;height:24px;background:white;border-radius:50%;transition:transform 0.3s}.toggle-switch.active .toggle-slider{transform:translateX(30px)}.section{display:none}.section.active{display:block}.input-group{margin-bottom:20px}label{display:block;margin-bottom:8px;color:#555;font-weight:600}textarea{width:100%;padding:12px;border:2px solid #e0e0e0;border-radius:10px;font-size:1em;font-family:inherit;resize:vertical;min-height:100px;transition:border-color 0.3s}textarea:focus{outline:none;border-color:#e94ca6}.emoji-picker{margin-bottom:20px}.emoji-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(45px,1fr));gap:8px;margin-bottom:15px}.emoji-btn{width:45px;height:45px;border:2px solid #e0e0e0;border-radius:8px;background:white;font-size:1.5em;cursor:pointer;transition:all 0.2s;display:flex;align-items:center;justify-content:center}.emoji-btn:hover{transform:scale(1.1);border-color:#e94ca6}.emoji-btn.selected{background:#e94ca6;border-color:#e94ca6}.letter-section{margin-top:20px}.letter-grid{display:grid;grid-template-columns:repeat(13,1fr);gap:5px}.letter-btn{width:100%;aspect-ratio:1;border:2px solid #e0e0e0;border-radius:6px;background:white;font-size:0.9em;font-weight:600;cursor:pointer;transition:all 0.2s}.letter-btn:hover{background:#f0f0f0;border-color:#e94ca6}.letter-btn.selected{background:#e94ca6;color:white;border-color:#e94ca6}.action-btn{width:100%;padding:14px;background:linear-gradient(135deg,#e94ca6 0%,#d946a3 100%);color:white;border:none;border-radius:10px;font-size:1em;font-weight:600;cursor:pointer;transition:all 0.3s;margin-top:20px}.action-btn:hover{transform:translateY(-2px);box-shadow:0 6px 20px rgba(233,76,166,0.4)}.result{margin-top:20px;padding:20px;background:white;border-radius:10px;border:2px dashed #e94ca6;display:none;text-align:center}.result.show{display:block}.result-emoji{font-size:3em;margin:10px 0}.result-text{font-size:1.2em;color:#333;word-break:break-all;margin:10px 0}.footer{margin-top:20px;text-align:center;color:white;font-size:0.9em}.footer a{color:white;text-decoration:underline}</style></head><body><div class="container"><h1>Hide a message in an emoji</h1><p class="description">This tool allows you to encode a hidden message into an emoji or alphabet letter. You can copy and paste text with a hidden message in it to decode the message.</p><div class="toggle-container"><span class="toggle-label">Decode</span><div class="toggle-switch active" id="modeToggle" onclick="toggleMode()"><div class="toggle-slider"></div></div><span class="toggle-label">Encode</span></div><div id="encode-section" class="section active"><div class="input-group"><label>Enter text to encode</label><textarea id="encode-text" placeholder="Type your secret message..."></textarea></div><div class="emoji-picker"><label>Pick an emoji</label><div class="emoji-grid"><button class="emoji-btn selected" data-emoji="😊" onclick="selectEmoji(this)">😊</button><button class="emoji-btn" data-emoji="🥳" onclick="selectEmoji(this)">🥳</button><button class="emoji-btn" data-emoji="😎" onclick="selectEmoji(this)">😎</button><button class="emoji-btn" data-emoji="🤩" onclick="selectEmoji(this)">🤩</button><button class="emoji-btn" data-emoji="😍" onclick="selectEmoji(this)">😍</button><button class="emoji-btn" data-emoji="🔥" onclick="selectEmoji(this)">🔥</button><button class="emoji-btn" data-emoji="💯" onclick="selectEmoji(this)">💯</button><button class="emoji-btn" data-emoji="👍" onclick="selectEmoji(this)">👍</button><button class="emoji-btn" data-emoji="🎉" onclick="selectEmoji(this)">🎉</button><button class="emoji-btn" data-emoji="❤️" onclick="selectEmoji(this)">❤️</button><button class="emoji-btn" data-emoji="🎨" onclick="selectEmoji(this)">🎨</button><button class="emoji-btn" data-emoji="🍕" onclick="selectEmoji(this)">🍕</button><button class="emoji-btn" data-emoji="🌙" onclick="selectEmoji(this)">🌙</button><button class="emoji-btn" data-emoji="🌈" onclick="selectEmoji(this)">🌈</button><button class="emoji-btn" data-emoji="😂" onclick="selectEmoji(this)">😂</button><button class="emoji-btn" data-emoji="🔮" onclick="selectEmoji(this)">🔮</button><button class="emoji-btn" data-emoji="💖" onclick="selectEmoji(this)">💖</button><button class="emoji-btn" data-emoji="🎭" onclick="selectEmoji(this)">🎭</button><button class="emoji-btn" data-emoji="👀" onclick="selectEmoji(this)">👀</button><button class="emoji-btn" data-emoji="🤓" onclick="selectEmoji(this)">🤓</button></div></div><div class="letter-section"><label>Or pick a standard alphabet letter</label><div class="letter-grid" id="letter-grid"></div></div><button class="action-btn" onclick="encodeMessage()">Encode</button><div id="encode-result" class="result"><div class="result-emoji" id="result-display"></div></div></div><div id="decode-section" class="section"><div class="input-group"><label>Paste encoded text</label><textarea id="decode-text" placeholder="Paste the encoded emoji or letter here..."></textarea></div><button class="action-btn" onclick="decodeMessage()">Decode</button><div id="decode-result" class="result"><div class="result-text" id="decoded-display"></div></div></div></div><div class="footer"><a href="https://github.com" target="_blank">Source on GitHub</a></div><script>let selectedEmoji='😊';let isEncodeMode=true;function initLetterGrid(){const grid=document.getElementById('letter-grid');for(let i=97;i<=122;i++){const letter=String.fromCharCode(i);const btn=document.createElement('button');btn.className='letter-btn';btn.textContent=letter;btn.onclick=()=>selectLetter(btn,letter);grid.appendChild(btn);}}function selectEmoji(btn){document.querySelectorAll('.emoji-btn').forEach(b=>b.classList.remove('selected'));document.querySelectorAll('.letter-btn').forEach(b=>b.classList.remove('selected'));btn.classList.add('selected');selectedEmoji=btn.dataset.emoji;}function selectLetter(btn,letter){document.querySelectorAll('.emoji-btn').forEach(b=>b.classList.remove('selected'));document.querySelectorAll('.letter-btn').forEach(b=>b.classList.remove('selected'));btn.classList.add('selected');selectedEmoji=letter;}function toggleMode(){isEncodeMode=!isEncodeMode;const toggle=document.getElementById('modeToggle');const encodeSection=document.getElementById('encode-section');const decodeSection=document.getElementById('decode-section');if(isEncodeMode){toggle.classList.add('active');encodeSection.classList.add('active');decodeSection.classList.remove('active');}else{toggle.classList.remove('active');encodeSection.classList.remove('active');decodeSection.classList.add('active');}}async function encodeMessage(){const text=document.getElementById('encode-text').value;if(!text){alert('Please enter a message to encode');return;}try{const response=await fetch('/api/encode',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({emoji:selectedEmoji,text:text})});const data=await response.json();document.getElementById('result-display').textContent=data.encoded;document.getElementById('encode-result').classList.add('show');navigator.clipboard.writeText(data.encoded);alert('Encoded! Copied to clipboard ✓');}catch(error){alert('Error: '+error.message);}}async function decodeMessage(){const text=document.getElementById('decode-text').value;if(!text){alert('Please paste an encoded message');return;}try{const response=await fetch('/api/decode',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({text:text})});const data=await response.json();const decoded=data.decoded||'(No hidden message found)';document.getElementById('decoded-display').textContent=decoded;document.getElementById('decode-result').classList.add('show');}catch(error){alert('Error: '+error.message);}}initLetterGrid();</script></body></html>`

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(htmlTemplate))
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/encode", encodeHandler)
	http.HandleFunc("/api/decode", decodeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
