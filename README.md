# 🦕 diny — git diff commit messages 

diny is a tiny dinosaur that writes your git commit messages for you.  
It looks at your staged changes and generates clear, conventional-friendly messages.

✅ No API key required — powered by my self hosted Ollama server  
🚀 Fast and reliable processing for all users  

---

## Why I’m building diny

I never liked using AI. I’m not a vibe coder, for me it kills the fun of programming. Still, everyone seems to use it and even judge others for not doing the same, while many couldn’t write a line of code without it.  

What I really struggle with are commit messages. I type fast, make typos, and usually leave junk like “ui update.” At the end of the day I want to log what I did for each client or project, but my commits are useless. I end up digging through git diffs just to remember what I worked on. With a memory like a goldfish, that’s exhausting.  

That’s why I built **diny**. It helps me write proper commit messages and has a **timeline** feature that pulls all commits for a day into a short summary. Perfect for updating clients, filling in my time tracker, or catching up if I missed a few days. For this kind of job AI actually makes sense — not to code for me, just to handle the boring parts I’d never do well myself.

## Features

- 🔍 Reads staged changes with `git diff`
- 🧹 Filters out noise (lockfiles, binaries, build artifacts)
- ⚡ Generates commit messages via Ollama
- 📝 Produces concise, consistent messages
- 🔄 Interactive workflow with multiple options:
  - **Commit** - Use the generated message
  - **Generate different message** - Get a completely new approach
  - **Refine message** - Provide custom feedback for targeted improvements
- 🧠 Smart regeneration that learns from previous attempts
- ✍️ Custom feedback system for precise message refinement
- 📊 Timeline analysis of commit history and message patterns

---

## Install

### Homebrew (macOS/Linux)

    brew install dinoDanic/tap/diny

### Windows

#### PowerShell (One-liner)
```powershell
$dest=Join-Path $env:LOCALAPPDATA 'diny\bin'; if(Test-Path $dest -PathType Leaf){throw "A FILE named '$dest' exists. Delete/rename it."}; New-Item -ItemType Directory -Path $dest -Force|Out-Null; $zip=Join-Path $env:TEMP 'diny.zip'; $tmp=Join-Path $env:TEMP ("diny_"+[guid]::NewGuid()); $arch=if($env:PROCESSOR_ARCHITECTURE -match 'ARM64'){'arm64|aarch64'}else{'amd64|x86_64|x64'}; $rel=Invoke-RestMethod "https://api.github.com/repos/dinoDanic/diny/releases/latest" -Headers @{ 'User-Agent'='PowerShell' }; $asset=$rel.assets|?{ $_.name -match "(?i)(windows|win).*($arch).*\.zip$"}|select -f 1; if(-not $asset){$asset=$rel.assets|?{ $_.name -match "(?i)(windows|win).*\.zip$"}|select -f 1}; if(-not $asset){throw "No Windows .zip asset found. Available:`n$($rel.assets.name -join "`n")"}; Invoke-WebRequest $asset.browser_download_url -OutFile $zip; Expand-Archive -Path $zip -DestinationPath $tmp -Force; Remove-Item $zip -Force; $exe=Get-ChildItem $tmp -Recurse -Filter "diny*.exe"|select -f 1; if(-not $exe){throw "Couldn't find diny.exe in the archive."}; $target=Join-Path $dest 'diny.exe'; if(Test-Path $target){Remove-Item $target -Force}; Move-Item $exe.FullName $target -Force; Get-ChildItem (Split-Path $exe.FullName) -Filter *.dll -ErrorAction SilentlyContinue | % { Copy-Item $_.FullName $dest -Force }; Remove-Item $tmp -Recurse -Force; if($env:PATH -notmatch [regex]::Escape($dest)){ $u=[Environment]::GetEnvironmentVariable('Path','User'); [Environment]::SetEnvironmentVariable('Path', ($u+";"+$dest).Trim(';'), 'User'); $env:PATH+=";$dest" }; & $target --help
```

#### Download Binary
Download the Windows zip from [GitHub Releases](https://github.com/dinoDanic/diny/releases) and extract `diny.exe`.

##### Add to PATH
Move `diny.exe` to a folder in your PATH (e.g. `C:\Windows\System32`)  
or create `C:\Users\<You>\bin`, add it to PATH via *System Properties → Environment Variables*, then restart the terminal.

### Download Binary (All Platforms)

Download pre-built binaries from [GitHub Releases](https://github.com/dinoDanic/diny/releases)

---

## Usage

Stage your changes, then run:

    git add -A
    diny commit

### Auto Command (Git Alias)

Set up a git alias that creates a `git auto` command for diny-generated commit messages.

```bash
diny auto          # Set up the git auto alias
diny auto remove   # Remove the git auto alias
```

After setup, you can run: 
```bash
git auto          # uses diny to generate commit message
```


### Configuration (Optional)

Configure diny settings for your project:

    diny init

---


## TODO

- [ ] 🌐 Language 
- [x] ⚙️ Per-user/project config  
- [x] 📦 Installation via popular package managers (Homebrew, Scoop, apt, etc.)  
- [x] ✍️ Interactive workflow with commit options
- [x] 🔄 Request alternative messages with smart regeneration
- [x] ✨ Custom feedback system for message refinement
- [ ] 🦥 Lazy git integration
