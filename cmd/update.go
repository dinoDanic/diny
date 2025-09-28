/*
Copyright © 2025 NAME HERE dino.danic@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/dinoDanic/diny/ui"
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
	ui.RenderTitle("Checking for diny updates...")

	checker := update.NewUpdateChecker(Version)
	latestVersion, err := checker.GetLatestVersion()
	if err != nil {
		ui.RenderError(fmt.Sprintf("Failed to check for updates: %v", err))
		os.Exit(1)
	}

	if !force && !checker.CompareVersions(Version, latestVersion) {
		ui.RenderSuccess(fmt.Sprintf("You're already on the latest version (%s)", Version))
		return
	}

	ui.RenderBox("Update Available", fmt.Sprintf("Updating from %s to %s...", Version, latestVersion))

	switch runtime.GOOS {
	case "darwin", "linux":
		updateUnixLike()
	case "windows":
		updateWindows(latestVersion)
	default:
		ui.RenderError(fmt.Sprintf("Unsupported operating system: %s", runtime.GOOS))
		showManualInstructions(latestVersion)
	}
}

func updateUnixLike() {
	if isHomebrewInstalled() {
		ui.RenderBox("Homebrew", "Updating via Homebrew...")
		if updateViaHomebrew() {
			return
		}
		ui.RenderError("Homebrew update failed")
	} else {
		showHomebrewInstallInstructions()
	}
}

func updateWindows(version string) {
	ui.RenderBox("Windows Update", "Installing/updating diny on Windows...")

	if runWindowsPowerShellInstaller() {
		return
	}

	ui.RenderError("PowerShell installation failed, showing manual instructions:")
	showWindowsManualInstructions(version)
}

func isHomebrewInstalled() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

func updateViaHomebrew() bool {
	ui.RenderBox("Homebrew Update", "Running brew update...")
	cmd := exec.Command("brew", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		ui.RenderError(fmt.Sprintf("brew update failed: %v", err))
		return false
	}

	ui.RenderBox("Homebrew Upgrade", "Running brew upgrade...")
	cmd = exec.Command("brew", "upgrade", "dinoDanic/tap/diny")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		ui.RenderError(fmt.Sprintf("brew upgrade failed: %v", err))
		return false
	}

	ui.RenderSuccess("Successfully updated via Homebrew")
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
		ui.RenderError(fmt.Sprintf("PowerShell execution failed: %v", err))
		return false
	}

	return true
}

func showHomebrewInstallInstructions() {
	ui.RenderWarning("Homebrew is not installed. Please install it to easily update diny")
}

func showWindowsManualInstructions(version string) {
	ui.RenderBox("Manual Windows Installation", fmt.Sprintf(`If automatic installation failed for version %s:

Visit: https://github.com/dinoDanic/diny`, version))
}

func showManualInstructions(version string) {
	ui.RenderBox("Manual Update Instructions", fmt.Sprintf(`For version %s:

macOS/Linux with Homebrew:
  brew update
  brew upgrade dinoDanic/tap/diny

Windows:
  Visit: https://github.com/dinoDanic/diny`, version))
}
