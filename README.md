# 🦕 diny — git diff commit messages 

diny is a tiny dinosaur that writes your git commit messages for you.  
It looks at your staged changes and generates clear, conventional-friendly messages using AI.

✅ No API key required — powered by my self hosted Ollama server  
🚀 Fast and reliable processing for all users  

---

## Features

- 🔍 Reads staged changes with `git diff --cached`
- 🧹 Filters out noise (lockfiles, binaries, build artifacts)
- ⚡ Generates commit messages via Ollama
- 📝 Produces concise, consistent messages
- 🔄 Interactive workflow with multiple options:
  - **Commit** - Use the generated message
  - **Generate different message** - Get a completely new approach
  - **Refine message** - Provide custom feedback for targeted improvements
  - **Exit** - Cancel the process
- 🧠 Smart regeneration that learns from previous attempts
- ✍️ Custom feedback system for precise message refinement

---

## Install

### Homebrew (macOS/Linux)

    brew install dinoDanic/tap/diny

### Download Binary

Download pre-built binaries from [GitHub Releases](https://github.com/dinoDanic/diny/releases)

---

## Usage

Stage your changes, then run:

    git add -A
    diny commit

---

## How it works

1. Collects a minimal version of your git diff (ignores lockfiles, builds, binaries).
2. Sends meaningful content to the Ollama model
3. Presents an interactive menu with multiple options:
   - **Commit** with the current message
   - **Generate different message** for alternative approaches
   - **Refine message** with your specific feedback
4. Smart regeneration tracks previous attempts to avoid repetition
5. Custom feedback allows you to guide the AI toward your preferred style

---

## TODO

- [ ] 🌐 Language 
- [x] ⚙️ Per-user/project config  
- [x] 📦 Installation via popular package managers (Homebrew, Scoop, apt, etc.)  
- [x] ✍️ Interactive workflow with commit options
- [x] 🔄 Request alternative messages with smart regeneration
- [x] ✨ Custom feedback system for message refinement
