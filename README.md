<div align="center">

<img src="https://diny-cli.vercel.app/diny-v2-200-light.png" alt="diny logo" width="200"/>

# diny

### from git diff to clean commits


diny is a tiny dinosaur that writes your git commit messages for you.  
It looks at your staged changes and generates clear, conventional-friendly messages.

🔓 **No API key required** • 🔌 Plug and play • ⚡ Fast and reliable • 🌟 Open source

</div>

<div align="center">
    <br />
<a href="https://peerlist.io/dinodanic/project/diny"><img src="https://diny-cli.vercel.app/peerlist-project-of-the-day.png" alt="diny showcase" width="200"></a>
    <br />
<img src="https://diny-cli.vercel.app/showcase/2.png" alt="diny showcase" width="700"/>
<img src="https://diny-cli.vercel.app/showcase/3.png" alt="diny showcase" width="700"/>
<img src="https://diny-cli.vercel.app/showcase/4.png" alt="diny showcase" width="700"/>
<img src="https://diny-cli.vercel.app/showcase/5.png" alt="diny showcase" width="700"/>
<img src="https://diny-cli.vercel.app/showcase/6.png" alt="diny showcase" width="700"/>

</div>

---

## Features

- 🔍 Reads staged changes with `git diff`
- 🧹 Filters out noise (lockfiles, binaries, build artifacts)
- ⚡ Generates commit messages via Ollama
- 📝 Produces concise, consistent messages
- 🔄 Interactive workflow with multiple options
- 🧠 Smart regeneration that learns from previous attempts
- ✍️ Custom feedback system for precise message refinement
- 🧷 Save to file (stash the generated message for later)
- 📝 Draft mode (prepare a commit message without committing)
- ✏️ Edit in your editor before committing
- 📊 Timeline analysis of commit history and message patterns
- 🎨 Customizable UI themes (10+ dark and light themes)


## Install

### macOS / Linux

```bash
brew install dinoDanic/tap/diny
```

### Windows

#### Scoop (Recommended)

```powershell
scoop bucket add dinodanic https://github.com/dinoDanic/scoop-bucket
scoop install diny
```

### Manual Installation

Download pre-built binaries from [GitHub Releases](https://github.com/dinoDanic/diny/releases)


## Usage

Stage your changes, then run:

    git add             # Stage files
    diny commit

### Auto Command (Git Alias)

Set up a git alias that creates a `git auto` command for diny-generated commit messages.

    diny auto          # Set up the git auto alias
    diny auto remove   # Remove the git auto alias

After setup, you can run:

    git auto           # uses diny to generate commit message

### LazyGit Integration

Integrate diny directly into [LazyGit](https://github.com/jesseduffield/lazygit) for seamless commit message generation from the LazyGit UI.

    diny link lazygit

This adds a custom command to LazyGit's configuration, allowing you to generate commit messages with diny without leaving LazyGit.


## Configuration

diny supports both global and local (per-repository) configuration with a clear precedence chain:

**Configuration Precedence:**
```
Environment Variables > Local Config > Global Config > Defaults
```

### Configuration Locations

**Global Config:**
- **Linux/Unix:** `~/.config/diny/config.json` (or `$XDG_CONFIG_HOME/diny/config.json`)
- **macOS:** `~/Library/Application Support/diny/config.json`
- **Windows:** `%APPDATA%\diny\config.json`

**Local Config:**
- `.git/diny-config.json` (in the git repository root)

### Environment Variables

You can override configuration with environment variables:

**Config Location:**
- `DINY_CONFIG_PATH` - Custom path for global config file (overrides platform defaults)

**Config Values:**
- `DINY_OLLAMA_URL` - Custom Ollama server URL (default: `http://127.0.0.1:11434`)
- `DINY_OLLAMA_MODEL` - Ollama model to use (default: `llama3.2`)
- `DINY_BACKEND_URL` - Cloud backend URL (default: `https://diny-cli.vercel.app`)

### Configuration Setup

Run the interactive configuration wizard:

```bash
diny init              # Choose between global or local configuration
```

View your current configuration:

```bash
diny config            # Show basic configuration
diny config --verbose  # Show effective configuration with sources
```

## Commands

diny comes with a handful of simple commands. Each one is designed to fit naturally into your git workflow:

    diny auto                                 # Set up a git alias so you can run `git auto`

    diny check-health                         # Check Ollama connection and model availability
    diny commit                               # Generate a commit message from your staged changes
    diny commit  --print                      # Print generated message to stdout only
    diny commit  --print | git commit -F -    # Generated message and commit

    diny config                               # Show your current diny configuration
    diny config --verbose                     # Show effective configuration with sources
    diny init                                 # Initialize diny with an interactive setup wizard
    diny link lazygit                         # Integrate diny with LazyGit
    diny theme                                # Select from 10+ dark and light UI themes
    diny theme list                           # List all available themes with previews
    diny timeline                             # Summarize and analyze your commit history
    diny update                               # Update diny to the latest version

## Update

### Built-in update command

```bash
diny update
```

### Manual update

**macOS/Linux:**
```bash
brew update && brew upgrade dinoDanic/tap/diny
```

**Windows (Scoop):**
```powershell
scoop update diny
```

## Why diny exists

I'm terrible at commit messages. I type fast, make mistakes, and my history ends up full of gems like "fix stuff" and "ui thing." Not great when you need to remember what you actually did.

I built **diny** because I wanted my commits to be useful without thinking about them. It reads your changes, writes a decent message, and has a timeline feature that summarizes your day. Super handy for client updates or just remembering what you worked on when your brain already moved on.

It's AI doing what it's actually good at—handling the repetitive stuff I'd rather not think about.


## TODO
[X] - Add lazygit integratio
