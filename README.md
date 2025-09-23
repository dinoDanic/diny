# 🦖 diny — AI commit messages 

diny is a tiny dinosaur that writes your git commit messages for you.  
It looks at your staged changes and generates clear, conventional-friendly messages using AI.

✅ No API key required — powered by my self-hosted Ollama server  
⚠️ Limited capacity right now (can be slow if many users) — will upgrade if needed  
🚧 Language, style, and custom prompts are not yet supported (coming soon - WIP)

---

## Features

- 🔍 Reads staged changes with `git diff --cached`
- 🧹 Filters out noise (lockfiles, binaries, build artifacts)
- 🤖 Generates commit messages via Ollama
- 📝 Produces concise, consistent messages

---

## Install

Requires Go (1.21+):

    go install github.com/dinoDanic/diny@latest

Make sure `~/go/bin` (or your $GOPATH/bin) is on your PATH.

---

## Usage

Stage your changes, then run:

    git add -A
    diny commit

---

## How it works

1. Collects a minimal version of your git diff (ignores lockfiles, builds, binaries).
2. Sends meaningful content to the Ollama model running on my server.
3. Suggests a commit message and asks for confirmation before committing.

---

## Roadmap

- 🌐 Language & style flags (--lang, --style)  
- ⏳ Streaming output (see the message as it’s generated)  
- ⚙️ Per-user/project config  
- 📦 Installation via popular package managers (Homebrew, Scoop, apt, etc.)  
- ✍️ Edit message before committing  
- 🔄 Request a new alternative message if not happy with the first one  
