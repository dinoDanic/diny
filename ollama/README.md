# Ollama Local Setup for Diny

This guide will help you set up Ollama locally to use with Diny instead of the cloud backend.

## Why Use Local Ollama?

- **Privacy**: Your code never leaves your machine
- **Speed**: No network latency for API calls
- **Offline**: Works without internet connection
- **Control**: Choose your own models and configurations

## Installation

### macOS

**Option 1: Direct Download** (Recommended)
```bash
# Download from ollama.com
curl -fsSL https://ollama.com/install.sh | sh
```

**Option 2: Homebrew**
```bash
brew install ollama
```

### Linux

**One-line install:**
```bash
curl -fsSL https://ollama.com/install.sh | sh
```

The script will:
- Install Ollama to `/usr/local/bin/ollama`
- Create a systemd service
- Start Ollama automatically

### Windows

**Option 1: Direct Download** (Recommended)
1. Visit [ollama.com/download](https://ollama.com/download)
2. Download `OllamaSetup.exe`
3. Run the installer
4. Follow the installation wizard

**Option 2: Scoop**
```powershell
scoop install ollama
```

### Docker

**Basic setup:**
```bash
docker pull ollama/ollama
docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
```

**With GPU support:**
```bash
docker run -d --gpus=all -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
```

## Configuration

### 1. Start Ollama Service

**macOS/Linux:**
```bash
ollama serve
```

**Windows:**
Ollama starts automatically after installation.

**Docker:**
Already running from the `docker run` command above.

### 2. Pull the Required Model

Diny uses `llama3.2` by default. Pull it with:

```bash
ollama pull llama3.2
```

**Other recommended models:**
```bash
# Smaller, faster model
ollama pull qwen2.5-coder:3b

# Larger, more capable model
ollama pull qwen2.5:7b-instruct
```

### 3. Verify Installation

Test that Ollama is running:

```bash
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2",
  "prompt": "Hello, world!",
  "stream": false
}'
```

You should see a JSON response with generated text.

## Configure Diny to Use Local Ollama

### During Initialization

When running `diny init`, answer "Yes" to:
```
Use local Ollama API?
Connect to local Ollama at http://localhost:11434
```

### Manual Configuration

Edit `.git/diny-config.json` in your repository:
```json
{
  "useConventional": true,
  "useEmoji": true,
  "tone": "professional",
  "length": "normal",
  "useLocalAPI": true
}
```

## Troubleshooting

### Ollama not responding

**Check if Ollama is running:**
```bash
# macOS/Linux
ps aux | grep ollama

# Windows
tasklist | findstr ollama
```

**Restart Ollama:**
```bash
# macOS/Linux
pkill ollama && ollama serve

# Windows
# Restart from Task Manager or system tray

# Docker
docker restart ollama
```

### Port already in use

If port 11434 is taken, you can change it:

```bash
# Set custom port
export OLLAMA_HOST=0.0.0.0:8080
ollama serve
```

Then update `ollama/ollama.go:14`:
```go
const server = "http://127.0.0.1:8080"
```

### Model not found

Make sure you've pulled the model:
```bash
ollama list  # Show installed models
ollama pull llama3.2  # Pull the model
```

### Out of memory

Try a smaller model:
```bash
ollama pull qwen2.5-coder:3b
```

Then update `ollama/ollama.go:19`:
```go
const model = "qwen2.5-coder:3b"
```

## System Requirements

### Minimum Recommended Specs

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| RAM | 8 GB | 16 GB+ |
| Storage | 10 GB free | 20 GB+ free |
| CPU | 4 cores | 8+ cores |
| GPU | None (CPU mode) | NVIDIA GPU with 8GB+ VRAM |

### Model Size Reference

| Model | Size | RAM Required | Speed |
|-------|------|--------------|-------|
| qwen2.5-coder:3b | 1.9 GB | 4-6 GB | Fast |
| llama3.2 | 2.0 GB | 4-6 GB | Fast |
| qwen2.5:7b-instruct | 4.7 GB | 8-12 GB | Medium |
| mistral:7b-instruct | 4.1 GB | 8-12 GB | Medium |

## Advanced Configuration

### Change Default Host/Port

Edit `ollama/ollama.go`:
```go
// For localhost (default)
const server = "http://127.0.0.1:11434"

// For network access
const server = "http://0.0.0.0:11434"

// For custom port
const server = "http://127.0.0.1:8080"
```

### Enable Network Access

To allow other machines to connect:

**Linux (systemd):**
```bash
sudo systemctl edit ollama.service
```

Add:
```
[Service]
Environment="OLLAMA_HOST=0.0.0.0"
```

Restart:
```bash
sudo systemctl restart ollama
```

## Useful Commands

```bash
# List installed models
ollama list

# Pull a model
ollama pull <model-name>

# Remove a model
ollama rm <model-name>

# Show model info
ollama show llama3.2

# Run a model interactively
ollama run llama3.2

# Check Ollama version
ollama --version
```

## Resources

- [Ollama GitHub](https://github.com/ollama/ollama)
- [Ollama Models Library](https://ollama.com/library)
- [Ollama API Documentation](https://github.com/ollama/ollama/blob/main/docs/api.md)

## Getting Help

If you encounter issues:
1. Check [Ollama GitHub Issues](https://github.com/ollama/ollama/issues)
2. Visit [r/ollama on Reddit](https://reddit.com/r/ollama)
3. Open an issue in the [Diny repository](https://github.com/dinoDanic/diny/issues)
