# Large Test File

This is a large test file designed to test git diff handling with substantial content that might challenge the AI system and token limits.

## Table of Contents
1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Configuration](#configuration)
4. [Usage](#usage)
5. [API Reference](#api-reference)
6. [Examples](#examples)
7. [Best Practices](#best-practices)
8. [Troubleshooting](#troubleshooting)
9. [Contributing](#contributing)
10. [License](#license)

## Introduction

Diny is an AI-powered Git commit message generator that analyzes your staged changes and creates meaningful, consistent commit messages. It integrates with Ollama to provide local AI processing, ensuring your code never leaves your development environment.

### Key Features
- 🤖 AI-powered commit message generation
- 🎨 Emoji support with gitmoji integration
- 📝 Conventional commits compliance
- ⚙️ Configurable tone and length
- 🏠 Local processing with Ollama
- 🚀 Fast and lightweight
- 🔧 Project-specific configuration

### Supported Models
- qwen2.5-coder:7b-instruct (recommended)
- deepseek-coder:6.7b-instruct
- codellama:7b-instruct
- mistral:7b-instruct
- And many more Ollama-compatible models

## Installation

### Prerequisites
- Git installed and configured
- Ollama running locally or on remote server
- Go 1.21 or higher (for building from source)

### Install from Release
```bash
# Download the latest release
wget https://github.com/dinoDanic/diny/releases/latest/download/diny-linux-amd64
chmod +x diny-linux-amd64
sudo mv diny-linux-amd64 /usr/local/bin/diny
```

### Build from Source
```bash
git clone https://github.com/dinoDanic/diny.git
cd diny
go build -o diny
sudo mv diny /usr/local/bin/
```

### Verify Installation
```bash
diny --help
diny --version
```

## Configuration

### Initial Setup
Run the interactive configuration wizard:
```bash
diny init
```

This will guide you through setting up:
- Emoji preferences
- Conventional commit format
- Tone (professional, casual, friendly)
- Message length (short, normal, long)

### Configuration File
Configuration is stored in `.git/diny-config.json`:
```json
{
  "use_emoji": true,
  "use_conventional": true,
  "tone": "casual",
  "length": "normal"
}
```

### Manual Configuration
You can manually edit the configuration file or use environment variables:
```bash
export DINY_USE_EMOJI=true
export DINY_USE_CONVENTIONAL=true
export DINY_TONE=casual
export DINY_LENGTH=normal
```

## Usage

### Basic Usage
```bash
# Stage your changes
git add .

# Generate and commit with AI message
diny commit
```

### Advanced Usage
```bash
# Use specific model
diny commit --model qwen2.5-coder:7b

# Force specific configuration
diny commit --emoji --conventional --tone professional

# Dry run (generate message without committing)
diny commit --dry-run

# Use custom Ollama endpoint
diny commit --endpoint http://my-server:11434
```

### Workflow Integration
```bash
# In your shell profile (.bashrc, .zshrc)
alias gc="diny commit"
alias gci="diny init"

# Git hooks (optional)
# .git/hooks/prepare-commit-msg
#!/bin/sh
if [ -z "$2" ]; then
    diny commit --dry-run > "$1"
fi
```

## API Reference

### Commands

#### `diny init`
Interactive configuration setup.

**Flags:**
- `--reset`: Reset to default configuration
- `--global`: Set global configuration (future feature)

#### `diny commit`
Generate and create commit with AI message.

**Flags:**
- `--model MODEL`: Specify Ollama model
- `--endpoint URL`: Custom Ollama endpoint
- `--dry-run`: Generate message without committing
- `--emoji`: Force emoji usage
- `--no-emoji`: Disable emoji
- `--conventional`: Force conventional format
- `--free-form`: Use free-form messages
- `--tone TONE`: Override tone (professional|casual|friendly)
- `--length LENGTH`: Override length (short|normal|long)

#### `diny config`
Manage configuration (future feature).

**Subcommands:**
- `diny config show`: Display current configuration
- `diny config set KEY VALUE`: Set configuration value
- `diny config reset`: Reset to defaults

### Configuration Schema

```typescript
interface UserConfig {
  use_emoji: boolean;        // Enable emoji prefixes
  use_conventional: boolean; // Use conventional commit format
  tone: 'professional' | 'casual' | 'friendly';
  length: 'short' | 'normal' | 'long';
}
```

### Exit Codes
- `0`: Success
- `1`: General error
- `2`: Configuration error
- `3`: Git error
- `4`: AI/Ollama error
- `5`: User cancellation

## Examples

### Example Outputs

**Short + Emoji + Conventional:**
```
✨ feat(auth): add OAuth login system
```

**Normal + No Emoji + Conventional:**
```
feat(api): implement user authentication endpoints

- Add POST /auth/login endpoint
- Add POST /auth/logout endpoint  
- Add middleware for token validation
- Update user model with auth fields
```

**Long + Emoji + Free-form:**
```
🔐 Implement comprehensive authentication system

This commit introduces a complete authentication system including:

- OAuth 2.0 integration with Google and GitHub providers
- JWT token generation and validation middleware
- Secure password hashing using bcrypt
- Session management with Redis store
- Rate limiting for login attempts
- Password reset functionality via email
- Two-factor authentication support
- Comprehensive test coverage for all auth flows

The implementation follows security best practices and includes
proper error handling, input validation, and audit logging.
```

### Diff Processing Examples

**Small diff (< 500 chars):**
- Processed quickly
- Focused commit message
- Single feature/fix description

**Medium diff (500-2000 chars):**
- May include bullet points
- Multiple file changes described
- Clear scope identification

**Large diff (> 2000 chars):**
- Warning displayed to user
- Longer processing time
- May be chunked or summarized
- Detailed body with impact analysis

## Best Practices

### For Teams
1. **Standardize configuration**: Share `.diny.json` template
2. **Use conventional commits**: Enable for changelog generation
3. **Set appropriate tone**: Professional for formal projects
4. **Review generated messages**: Always verify before committing

### For Solo Projects
1. **Customize to your style**: Use casual tone and emojis
2. **Adjust length based on project**: Short for small projects
3. **Experiment with models**: Find what works best for your code

### Performance Tips
1. **Use code-specific models**: `qwen2.5-coder` over general models
2. **Keep diffs focused**: Commit logical units of change
3. **Exclude generated files**: Use `.gitignore` for build artifacts
4. **Monitor token usage**: Large diffs consume more resources

### Security Considerations
1. **Local processing**: Code never leaves your environment
2. **Network usage**: Only communicates with your Ollama instance
3. **Configuration storage**: Stored locally in `.git/` folder
4. **No telemetry**: No data collection or external calls

## Troubleshooting

### Common Issues

**"Ollama not found" error:**
```bash
# Check if Ollama is running
curl http://localhost:11434/api/tags

# Start Ollama if needed
ollama serve

# Check if model is available
ollama list
```

**"No staged changes" error:**
```bash
# Stage your changes first
git add .

# Or stage specific files
git add src/main.go config/app.json
```

**"Configuration not found" error:**
```bash
# Re-run initialization
diny init

# Or check if you're in a git repository
git status
```

**"Diff too large" warning:**
- Break changes into smaller commits
- Exclude large generated files
- Use `--max-tokens` flag (future feature)

### Model Performance Issues

**Slow generation:**
- Switch to smaller model: `deepseek-coder:6.7b`
- Use faster hardware or server
- Reduce diff size by committing more frequently

**Poor message quality:**
- Try code-specific models
- Adjust tone and length settings
- Review and edit generated messages
- Report issues with examples

### Debug Mode
```bash
# Enable verbose logging (future feature)
diny commit --debug

# Show diff processing details
diny commit --verbose

# Test configuration
diny config validate
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup
```bash
git clone https://github.com/dinoDanic/diny.git
cd diny
go mod download
go build -o diny
./diny init
```

### Running Tests
```bash
go test ./...
go test -race ./...
go test -cover ./...
```

### Code Style
- Use `gofmt` and `goimports`
- Follow Go conventions
- Add tests for new features
- Update documentation

### Submitting Changes
1. Fork the repository
2. Create feature branch
3. Make your changes
4. Add tests and documentation
5. Submit pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

This large test file should provide substantial content for testing git diff processing, AI token limits, and commit message generation with various scenarios and edge cases.