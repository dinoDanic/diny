# Ollama Local Setup for Diny

This guide will help you set up Ollama locally to use with Diny instead of the cloud backend.

## Why Use Local Ollama?

- **Privacy**: Your code never leaves your machine
- **Offline**: Works without internet connection
- **Control**: Choose your own models and configurations

## ⚠️ Performance Expectations

### Context Size vs Speed Trade-offs (CPU-only, 4-core)

| Context | Speed Impact | 4k diff | 15k diff | Use Case |
|---------|--------------|---------|----------|----------|
| 8k | Baseline | ~30-60s | ~45-90s | Daily commits |
| 16k | 2x slower | ~60-120s | ~90-180s | Large features |
| 32k | 4x slower | ~120-240s | ~180-360s | Major refactors |

**Note**: Context size affects speed even for small diffs because it pre-allocates memory. Higher context = slower but handles larger diffs.

### GPU Performance

**With NVIDIA GPU:**
- All context sizes: 2-10 seconds typically
- Minimal speed difference between 8k/16k/32k

**Recommendation**: If you don't have a GPU and need fast responses, use the default cloud backend instead.

## Installation

### Linux (Recommended)

```bash
curl -fsSL https://ollama.com/install.sh | sh
```

The script will:
- Install Ollama to `/usr/local/bin/ollama`
- Create a systemd service
- Start Ollama automatically

### macOS

**Direct Download:**
```bash
curl -fsSL https://ollama.com/install.sh | sh
```

**Homebrew:**
```bash
brew install ollama
```

### Windows

1. Visit [ollama.com/download](https://ollama.com/download)
2. Download `OllamaSetup.exe`
3. Run the installer
4. Ollama starts automatically after installation

### Docker (NOT Recommended for CPU-only systems)

Docker adds overhead and is significantly slower without GPU passthrough.

**CPU-only (slow):**
```bash
docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
```

**With GPU (requires nvidia-docker):**
```bash
docker run -d --gpus=all -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
```

## Configuration

### 1. Start Ollama Service

**Linux:**
```bash
# Systemd service (auto-starts after install)
sudo systemctl status ollama

# Or run manually
ollama serve
```

**macOS:**
```bash
ollama serve
```

**Windows:**
Ollama starts automatically after installation.

### 2. Pull and Configure Model

```bash
# Pull base model
ollama pull llama3.2

# Create 8k context version (handles larger git diffs)
ollama create llama3.2-8k -f - << 'EOF'
FROM llama3.2
PARAMETER num_ctx 8192
EOF
```

**Why 8k context?**
- Default: 4096 tokens (fails on diffs > ~3500 tokens)
- With 8192: Handles diffs up to ~7000 tokens
- Typical multi-file changes: 2000-5000 tokens

### 3. Configure Diny

Run `diny init` and select:
- "Use local Ollama API?" → Yes
- "Ollama model?" → llama3.2-8k

Or edit `.git/diny-config.json`:
```json
{
  "useConventional": true,
  "useEmoji": true,
  "tone": "professional",
  "length": "normal",
  "useLocalAPI": true,
  "ollamaModel": "llama3.2-8k"
}
```

### 4. Verify Setup

```bash
diny check-health
```

Should show:
```
✓ Connected to Ollama successfully
✓ Model 'llama3.2-8k' is available and ready to use
```

## Alternative Models

### For CPU-only systems (faster but less capable)

```bash
# Smaller, faster model (~2x speed of llama3.2)
ollama pull qwen2.5-coder:3b

# Create 8k context version
ollama create qwen2.5-coder:3b-8k -f - << 'EOF'
FROM qwen2.5-coder:3b
PARAMETER num_ctx 8192
EOF
```

Then run `diny init` and select `qwen2.5-coder:3b-8k` as the model.

### Model Comparison

| Model | Size | RAM | CPU Speed* | GPU Speed | Quality |
|-------|------|-----|-----------|-----------|---------|
| qwen2.5-coder:3b-8k | 1.9 GB | 4-6 GB | ~15-30s | ~2-5s | Good |
| llama3.2-8k | 2.0 GB | 4-6 GB | ~30-60s | ~3-8s | Better |
| qwen2.5:7b-8k | 4.7 GB | 8-12 GB | ~60-120s | ~5-15s | Best |

\* Typical 4000 token diff on 4-core CPU

## Troubleshooting

### Slow performance / Long wait times

**Expected on CPU-only systems.** Options:

1. **Use smaller model**: Switch to `qwen2.5-coder:3b-8k`
2. **Use cloud backend**: Run `diny init` and select "No" for local API
3. **Smaller commits**: Break large changes into smaller commits
4. **Add GPU**: Install NVIDIA GPU for 10-20x speedup

### Connection refused

```bash
# Check if Ollama is running
ps aux | grep ollama

# Start manually
ollama serve

# Or restart systemd service (Linux)
sudo systemctl restart ollama
```

### Model not found

```bash
# List installed models
ollama list

# Pull missing model
ollama pull llama3.2

# Recreate 8k version
ollama create llama3.2-8k -f - << 'EOF'
FROM llama3.2
PARAMETER num_ctx 8192
EOF
```

### Diff too large / Context exceeded

Your diff exceeds 8192 tokens. Options:

1. **Commit in smaller chunks**: Stage fewer files
2. **Increase context** (slower):
```bash
ollama create llama3.2-16k -f - << 'EOF'
FROM llama3.2
PARAMETER num_ctx 16384
EOF
```
Then update config to use `llama3.2-16k`

3. **Use cloud backend**: No token limits

### Out of memory

```bash
# Switch to smaller model
ollama pull qwen2.5-coder:3b

ollama create qwen2.5-coder:3b-8k -f - << 'EOF'
FROM qwen2.5-coder:3b
PARAMETER num_ctx 8192
EOF
```

Then run `diny init` and select `qwen2.5-coder:3b-8k`.

## System Requirements

### Minimum (CPU-only)
- **CPU**: 4+ cores (8+ recommended)
- **RAM**: 8 GB (16 GB recommended)
- **Storage**: 10 GB free
- **Note**: Generation will be slow (30-90s for large diffs)

### Recommended (GPU)
- **GPU**: NVIDIA GPU with 8GB+ VRAM
- **CPU**: 4+ cores
- **RAM**: 8 GB
- **Storage**: 10 GB free
- **Note**: Fast generation (2-10s for large diffs)

## Environment Variables

```bash
# Custom Ollama URL
export DINY_OLLAMA_URL="http://192.168.1.100:11434"

# Custom model
export DINY_OLLAMA_MODEL="qwen2.5-coder:3b-8k"

# Custom Ollama port (for ollama serve)
export OLLAMA_HOST="0.0.0.0:8080"
```

## Useful Commands

```bash
# Check Ollama status
diny check-health

# List installed models
ollama list

# Show model details
ollama show llama3.2-8k

# Remove unused models
ollama rm llama3.2

# Test model interactively
ollama run llama3.2-8k

# Check Ollama version
ollama --version
```

## Advanced Configuration

### Dynamic Context (Alternative Approach)

Instead of creating separate models for different context sizes, you can use the base model with runtime context override:

**Via CLI:**
```bash
# Run base model with 8k context
ollama run llama3.2 --num-ctx 8192

# Or 16k context
ollama run llama3.2 --num-ctx 16384
```

**Via API:**
```bash
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2",
  "prompt": "test",
  "options": {
    "num_ctx": 8192
  }
}'
```

**When to use this:**
- You want to avoid creating multiple model variants
- You need different context sizes for different use cases
- You're comfortable modifying Diny's code

**When to use separate models (recommended for most users):**
- Simple setup, no code changes needed
- Clear model naming (`llama3.2-8k` vs `llama3.2-16k`)
- Works out-of-the-box with Diny

### Modifying Diny for Dynamic Context

If you want Diny to support dynamic context configuration:

1. **Add config field** in `config/config.go`:
```go
type UserConfig struct {
    // ... existing fields ...
    OllamaContext int `json:"ollamaContext,omitempty"`
}
```

2. **Update API request** in `ollama/ollama.go`:
```go
type GenerateRequest struct {
    Model   string                 `json:"model"`
    Prompt  string                 `json:"prompt"`
    Stream  bool                   `json:"stream"`
    Options map[string]interface{} `json:"options,omitempty"`
}

// In Main() and MainStream():
req := GenerateRequest{
    Model:  apiConfig.Model,
    Prompt: prompt,
    Stream: false,
}

// Add context if configured
if userConfig.OllamaContext > 0 {
    req.Options = map[string]interface{}{
        "num_ctx": userConfig.OllamaContext,
    }
}
```

3. **Update init prompts** to ask for context size

This allows configuration like:
```json
{
  "useLocalAPI": true,
  "ollamaModel": "llama3.2",
  "ollamaContext": 8192
}
```

## Resources

- [Ollama GitHub](https://github.com/ollama/ollama)
- [Ollama Models Library](https://ollama.com/library)
- [Ollama API Documentation](https://github.com/ollama/ollama/blob/main/docs/api.md)
- [Modelfile Reference](https://github.com/ollama/ollama/blob/main/docs/modelfile.md)

## Getting Help

If you encounter issues:
1. Run `diny check-health` to diagnose
2. Check [Ollama GitHub Issues](https://github.com/ollama/ollama/issues)
3. Open an issue in the [Diny repository](https://github.com/dinoDanic/diny/issues)
