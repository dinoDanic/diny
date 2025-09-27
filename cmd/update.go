/*
Copyright © 2025 NAME HERE dino.danic@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/dinoDanic/diny/update"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update diny to the latest version",
	Long: `Update diny to the latest version.

This command will:
- Check for the latest version on GitHub
- macOS/Linux: Update via Homebrew or show Homebrew installation instructions  
- Windows: Run PowerShell installer automatically

Examples:
  diny update
  diny update --force    # Force update even if already latest
`,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		runUpdate(force)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolP("force", "f", false, "Force update even if already on latest version")
}

func runUpdate(force bool) {
	fmt.Println("🦕 Checking for diny updates...")

	checker := update.NewUpdateChecker(Version)
	latestVersion, err := checker.GetLatestVersion()
	if err != nil {
		fmt.Printf("❌ Failed to check for updates: %v\n", err)
		os.Exit(1)
	}

	if !force && !checker.CompareVersions(Version, latestVersion) {
		fmt.Printf("✅ You're already on the latest version (%s)\n", Version)
		return
	}

	fmt.Printf("📦 Updating from %s to %s...\n", Version, latestVersion)

	switch runtime.GOOS {
	case "darwin", "linux":
		updateUnixLike()
	case "windows":
		updateWindows(latestVersion)
	default:
		fmt.Printf("❌ Unsupported operating system: %s\n", runtime.GOOS)
		showManualInstructions(latestVersion)
	}
}

func updateUnixLike() {
	if isHomebrewInstalled() {
		fmt.Println("🍺 Updating via Homebrew...")
		if updateViaHomebrew() {
			return
		}
		fmt.Println("❌ Homebrew update failed")
	} else {
		showHomebrewInstallInstructions()
	}
}

func updateWindows(version string) {
	fmt.Println("💻 Installing/updating diny on Windows...")

	if runWindowsPowerShellInstaller() {
		return
	}

	fmt.Println("❌ PowerShell installation failed, showing manual instructions:")
	showWindowsManualInstructions(version)
}

func isHomebrewInstalled() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

func updateViaHomebrew() bool {
	cmd := exec.Command("brew", "update")
	if err := cmd.Run(); err != nil {
		return false
	}

	cmd = exec.Command("brew", "upgrade", "dinoDanic/tap/diny")
	if err := cmd.Run(); err != nil {
		return false
	}

	fmt.Println("🦕 Successfully updated via Homebrew")
	return true
}

func runWindowsPowerShellInstaller() bool {
	script := `
$ErrorActionPreference = "Stop"
try {
    Write-Host "📥 Downloading and installing diny..."
    
    $dest = Join-Path $env:LOCALAPPDATA 'diny\bin'
    if (Test-Path $dest -PathType Leaf) {
        throw "A FILE named '$dest' exists. Delete/rename it."
    }
    
    New-Item -ItemType Directory -Path $dest -Force | Out-Null
    $zip = Join-Path $env:TEMP 'diny.zip'
    $tmp = Join-Path $env:TEMP ("diny_" + [guid]::NewGuid())
    
    # Get latest release info
    $rel = Invoke-RestMethod "https://api.github.com/repos/dinoDanic/diny/releases/latest" -Headers @{ 'User-Agent' = 'PowerShell' }
    
    # Find Windows asset
    $asset = $rel.assets | Where-Object { $_.name -match "(?i)windows.*\.zip$" } | Select-Object -First 1
    if (-not $asset) {
        throw "No Windows .zip asset found in latest release"
    }
    
    Write-Host "📦 Downloading $($asset.name)..."
    Invoke-WebRequest $asset.browser_download_url -OutFile $zip
    
    Write-Host "📂 Extracting archive..."
    Expand-Archive -Path $zip -DestinationPath $tmp -Force
    Remove-Item $zip -Force
    
    $exe = Get-ChildItem $tmp -Recurse -Filter "diny*.exe" | Select-Object -First 1
    if (-not $exe) {
        throw "Couldn't find diny.exe in the archive"
    }
    
    $target = Join-Path $dest 'diny.exe'
    if (Test-Path $target) {
        Remove-Item $target -Force
    }
    
    Move-Item $exe.FullName $target -Force
    Remove-Item $tmp -Recurse -Force
    
    # Update PATH if needed
    if ($env:PATH -notmatch [regex]::Escape($dest)) {
        $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
        $newPath = ($userPath + ";" + $dest).Trim(';')
        [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
        $env:PATH += ";$dest"
        Write-Host "✅ Added $dest to PATH"
    }
    
    Write-Host "✅ Successfully installed diny to $target"
    Write-Host "🔄 You may need to restart your terminal for PATH changes to take effect"
    
    # Test the installation
    & $target --version
    exit 0
    
} catch {
    Write-Host "❌ Installation failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
`

	// Run PowerShell with the script
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("❌ PowerShell execution failed: %v\n", err)
		return false
	}

	return true
}

func showHomebrewInstallInstructions() {
	fmt.Printf(`
🍺 Homebrew is not installed. Please install it to easily update diny:
`)
}

func showWindowsManualInstructions(version string) {
	fmt.Print(`
💻 Manual Windows Installation:

📋 If automatic installation failed,
  Visit: https://github.com/dinoDanic/diny

`, version)
}

func showManualInstructions(version string) {
	fmt.Print(`
📋 Manual Update Instructions:

macOS/Linux with Homebrew:
  brew update
  brew upgrade dinoDanic/tap/diny

Windows:
  Visit: https://github.com/dinoDanic/diny

`, version)
}
