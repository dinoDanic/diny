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
- ⚡ Generates commit messages via AI
- 📝 Produces concise, consistent messages
- 🔄 Interactive workflow with multiple options
- 🧠 Smart regeneration that learns from previous attempts
- ✍️ Custom feedback system for precise message refinement
- 🧷 Save to file (stash the generated message for later)
- 📝 Draft mode (prepare a commit message without committing)
- ✏️ Edit in your editor before committing
- 📊 Timeline analysis of commit history and message patterns
- 📋 Interactive changelog viewer for GitHub releases
- 🎨 Customizable UI themes (10+ dark and light themes)
- ⚙️ Three-tier config system (global, project-shared, project-private)


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


## Commands

diny comes with a handful of simple commands. Each one is designed to fit naturally into your git workflow:

    diny auto                                 # Set up a git alias so you can run `git auto`

    diny commit                               # Generate a commit message from your staged changes
    diny commit  --no-verify                  # Commit without running git hooks
    diny commit  --print                      # Print generated message to stdout only
    diny commit  --print | git commit -F -    # Pipe generated message to git commit
    diny commit  --length short               # Force short commit message length (-l short)
    diny commit  --length normal              # Force normal commit message length (-l normal)
    diny commit  --length long                # Force long commit message length (-l long)
    diny commit  --custom "mention ticket"    # One-off custom instruction for the AI (-c "mention ticket")

    diny yolo                                 # Stage all changes, generate commit, and push (fully automated)

    diny changelog                            # View and interact with GitHub release changelogs
    diny config                               # Open config file in your editor
    diny link lazygit                         # Integrate diny with LazyGit
    diny theme                                # List all available UI themes
    diny timeline                             # Summarize and analyze your commit history
    diny update                               # Update diny to the latest version

## Configuration

diny supports a three-tier configuration system that allows you to customize settings globally, per-project (shared with your team), or per-project (private to you).

### Configuration Levels

When you run `diny config` in a git repository, you can choose between three config levels:

| Level | Location | Description | Committed? |
|-------|----------|-------------|------------|
| **Global** | `~/.config/diny/config.yaml` | Your default settings for all projects | No |
| **Project (Versioned)** | `.diny.yaml` | Team-shared settings for this project | Yes ✓ |
| **Project (Local)** | `<gitdir>/diny/config.yaml` | Your personal overrides for this project | No |

**Priority:** Local > Versioned > Global (higher priority overrides lower)

**Outside a git repository:** Only global config is available.

### When to Use Each Level

- **Global**: Your personal defaults across all projects (e.g., preferred theme, tone)
- **Versioned**: Team standards everyone should follow (e.g., conventional commits, emoji usage)
- **Local**: Your personal preferences on a specific project (e.g., different tone or length)

Project configs only need to specify settings you want to override—they automatically inherit from lower levels.

### Options

| Option | Description | Values |
|--------|-------------|--------|
| `theme` | UI color theme | catppuccin, tokyo, nord, dracula, gruvbox-dark, onedark, monokai, solarized-dark, solarized-light, github-light, etc. |
| `commit.conventional` | Use conventional commit format | `true` / `false` |
| `commit.conventional_format` | Commit types to use | `['feat', 'fix', 'docs', 'chore', ...]` |
| `commit.emoji` | Add emoji prefix | `true` / `false` |
| `commit.emoji_map` | Emoji for each type | `feat: 🚀, fix: 🐛, ...` |
| `commit.tone` | Message tone | `professional` / `casual` / `friendly` |
| `commit.length` | Message length | `short` / `normal` / `long` |
| `commit.custom_instructions` | Custom AI guidance | Any text, e.g. "Always mention ticket number" |
| `commit.hash_after_commit` | Show and copy commit hash | `true` / `false` |

### Example Config

```yaml
theme: tokyo

commit:
  conventional: true
  conventional_format: ['feat', 'fix', 'docs', 'chore']
  emoji: true
  emoji_map:
    feat: 🚀
    fix: 🐛
    docs: 📚
    chore: 🔧
  tone: casual
  length: short
  custom_instructions: "Include JIRA ticket from branch name"
  hash_after_commit: true
```

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

