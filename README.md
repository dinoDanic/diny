<div align="center">

<img src="https://diny.run/diny-v2-200-light.png" alt="diny logo" width="200"/>

# diny

**Free AI git commit messages. No API key. No signup.**

</div>

Open-source CLI that turns staged diffs into clean commits in one command.

---

## Features

- Interactive TUI for commit, config, changelog, timeline, and yolo
- Reads staged changes with `git diff` and filters out noise (lockfiles, binaries, build artifacts)
- Generates 3 message variants; inline editing or open in `$EDITOR`
- File picker to stage/unstage without leaving diny
- Timeline analysis with date presets or custom ranges
- AI-powered changelog generation between tags or commits
- Three-tier config system (global, project-shared, project-private)
- 10+ dark and light themes
- No API key required

## Install

### macOS / Linux

```bash
brew install dinoDanic/tap/diny
```

### Windows (Scoop)

```powershell
scoop bucket add dinodanic https://github.com/dinoDanic/scoop-bucket
scoop install diny
```

### Manual

Download pre-built binaries from [GitHub Releases](https://github.com/dinoDanic/diny/releases).

## Usage

Stage your changes, then run:

```bash
git add .
diny commit
```

This launches the interactive TUI — generate, review, and commit without leaving the terminal.

## Commands

| Command | Description |
|---------|-------------|
| `diny commit` | Launch the interactive TUI |
| `diny yolo` | Stage all changes, generate a commit, and push |
| `diny changelog` | Generate an AI-powered changelog between tags or commits |
| `diny timeline` | Summarize and analyze your commit history |
| `diny config` | Interactive TUI config editor |
| `diny theme` | List available UI themes |
| `diny auto` | Set up a `git auto` alias |
| `diny link lazygit` | Integrate diny with LazyGit |
| `diny update` | Update diny to the latest version |

## Integrations

### Git alias (`git auto`)

```bash
diny auto          # install
diny auto remove   # uninstall
```

Then use `git auto` anywhere you'd use `git commit`.

### LazyGit

```bash
diny link lazygit
```

Adds a custom command to [LazyGit](https://github.com/jesseduffield/lazygit) so you can generate commit messages from its UI.

## Configuration

Run `diny config` to open a TUI for browsing and editing all settings.

diny supports a three-tier configuration system:

| Level | Location | Description | Committed |
|-------|----------|-------------|-----------|
| Global | `~/.config/diny/config.yaml` | Your defaults for all projects | No |
| Project (versioned) | `.diny.yaml` | Team-shared settings for this project | Yes |
| Project (local) | `<gitdir>/diny/config.yaml` | Your personal per-project overrides | No |

**Priority:** local > versioned > global. Project configs only need to specify the values they override.

Outside a git repository, only global config is available.

### Options

| Option | Description | Values |
|--------|-------------|--------|
| `theme` | UI color theme | see [Themes](#themes) |
| `commit.conventional` | Use conventional commit format | `true` / `false` |
| `commit.emoji` | Prefix commits with emoji | `true` / `false` |
| `commit.tone` | Message tone | `professional` / `casual` / `friendly` |
| `commit.length` | Message length | `short` / `normal` / `long` |
| `commit.custom_instructions` | Extra guidance for the AI | free text |
| `commit.hash_after_commit` | Show and copy commit hash after committing | `true` / `false` |

### Themes

- **Dark:** `catppuccin`, `tokyo`, `nord`, `dracula`, `gruvbox-dark`, `onedark`, `monokai`, `solarized-dark`, `everforest-dark`, `flexoki-dark`
- **Light:** `solarized-light`, `github-light`, `gruvbox-light`, `flexoki-light`

Run `diny theme` to preview them all.

### Example

```yaml
theme: tokyo

commit:
  conventional: true
  emoji: true
  tone: casual
  length: short
  custom_instructions: "Include JIRA ticket from branch name"
  hash_after_commit: true
```

## Update

```bash
diny update                                     # built-in updater
brew update && brew upgrade dinoDanic/tap/diny  # macOS / Linux
scoop update diny                               # Windows
```

## Why diny

I'm terrible at commit messages — I type fast, make mistakes, and my history ends up full of "fix stuff" and "ui thing." diny reads your changes, writes a decent message, and summarizes your day when you need to remember what you worked on. It's AI doing what it's actually good at: the repetitive stuff.

## License

[MIT](LICENSE)
