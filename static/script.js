let selectedEmoji = '😊';
let isEncodeMode = true;

const emojis = ['😊', '🥳', '😎', '🤩', '😍', '🔥', '💯', '👍', '🎉', '❤️', '🎨', '🍕', '🌙', '🌈', '😂', '🔮', '💖', '🎭', '👀', '🤓'];

function initEmojiGrid() {
    const grid = document.getElementById('emoji-grid');
    emojis.forEach((emoji, index) => {
        const btn = document.createElement('button');
        btn.className = 'emoji-btn' + (index === 0 ? ' selected' : '');
        btn.textContent = emoji;
        btn.onclick = () => selectEmoji(btn, emoji);
        grid.appendChild(btn);
    });
}

function initLetterGrid() {
    const grid = document.getElementById('letter-grid');
    for (let i = 97; i <= 122; i++) {
        const letter = String.fromCharCode(i);
        const btn = document.createElement('button');
        btn.className = 'letter-btn';
        btn.textContent = letter;
        btn.onclick = () => selectLetter(btn, letter);
        grid.appendChild(btn);
    }
}

function selectEmoji(btn, emoji) {
    document.querySelectorAll('.emoji-btn').forEach(b => b.classList.remove('selected'));
    document.querySelectorAll('.letter-btn').forEach(b => b.classList.remove('selected'));
    btn.classList.add('selected');
    selectedEmoji = emoji;
}

function selectLetter(btn, letter) {
    document.querySelectorAll('.emoji-btn').forEach(b => b.classList.remove('selected'));
    document.querySelectorAll('.letter-btn').forEach(b => b.classList.remove('selected'));
    btn.classList.add('selected');
    selectedEmoji = letter;
}

function toggleMode() {
    isEncodeMode = !isEncodeMode;
    const toggle = document.getElementById('modeToggle');
    const encodeSection = document.getElementById('encode-section');
    const decodeSection = document.getElementById('decode-section');
    
    if (isEncodeMode) {
        toggle.classList.add('active');
        encodeSection.classList.add('active');
        decodeSection.classList.remove('active');
    } else {
        toggle.classList.remove('active');
        encodeSection.classList.remove('active');
        decodeSection.classList.add('active');
    }
}

async function encodeMessage() {
    const text = document.getElementById('encode-text').value;
    if (!text) {
        alert('Please enter a message to encode');
        return;
    }
    
    try {
        const response = await fetch('/api/encode', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({emoji: selectedEmoji, text: text})
        });
        
        const data = await response.json();
        document.getElementById('result-display').textContent = data.encoded;
        document.getElementById('encode-result').classList.add('show');
        
        navigator.clipboard.writeText(data.encoded);
        alert('Encoded! Copied to clipboard ✓');
    } catch (error) {
        alert('Error: ' + error.message);
    }
}

async function decodeMessage() {
    const text = document.getElementById('decode-text').value;
    if (!text) {
        alert('Please paste an encoded message');
        return;
    }
    
    try {
        const response = await fetch('/api/decode', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({text: text})
        });
        
        const data = await response.json();
        const decoded = data.decoded || '(No hidden message found)';
        document.getElementById('decoded-display').textContent = decoded;
        document.getElementById('decode-result').classList.add('show');
    } catch (error) {
        alert('Error: ' + error.message);
    }
}

initEmojiGrid();
initLetterGrid();   