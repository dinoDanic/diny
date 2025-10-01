<div align="center">

<img src="https://diny-cli.vercel.app/diny-circle-200-5.png" alt="diny logo" width="200"/>

# diny

### from git diff to clean commits

diny is a tiny dinosaur that writes your git commit messages for you.  
It looks at your staged changes and generates clear, conventional-friendly messages.

🔓 **No API key required** • 🔌 Plug and play • ⚡ Fast and reliable • 🌟 Open source

</div>

<div align="center">

![diny demo](https://diny-cli.vercel.app/demo.gif)

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

#### Winget

```powershell
winget install dinoDanic.diny
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
    diny commit  --print                      # Print generated message to stdout only
    diny commit  --print | git commit -F -    # Generated message and commit

    diny config                               # Show your current diny configuration
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

**Windows (Winget):**
```powershell
winget upgrade dinoDanic.diny
```

## Why I’m building diny

I never liked using AI. I’m not a vibe coder, for me it kills the fun of programming. Still, everyone seems to use it and even judge others for not doing the same, while many couldn’t write a line of code without it.  

What I really struggle with are commit messages. I type fast, make typos, and usually leave junk like “ui update.” At the end of the day I want to log what I did for each client or project, but my commits are useless. I end up digging through git diffs just to remember what I worked on. With a memory like a goldfish, that’s exhausting.  

That’s why I built **diny**. It helps me write proper commit messages and has a **timeline** feature that pulls all commits for a day into a short summary. Perfect for updating clients, filling in my time tracker, or catching up if I missed a few days. For this kind of job AI actually makes sense — not to code for me, just to handle the boring parts I’d never do well myself.


## TODO
[X] - Add lazygit integratio
