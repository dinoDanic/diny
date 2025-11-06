package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/ui"
	"github.com/spf13/cobra"
)

var autoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Set up git auto alias for diny commit messages",
	Long: `Set up a git alias that creates 'git auto' command for diny-generated commit messages.

After setting up the alias:
  git auto             -> uses diny to generate commit message

Examples:
  diny auto            # Set up the git auto alias
  diny auto remove     # Remove the git auto alias`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 && args[0] == "remove" {
			removeGitAlias()
		} else {
			setupGitAlias()
		}
	},
}

func setupGitAlias() {
	dinyPath, err := getDinyPath()
	if err != nil {
		ui.Box(ui.BoxOptions{
			Message: fmt.Sprintf("Error finding diny executable: %v", err),
			Variant: ui.Error,
		})
		os.Exit(1)
	}

	aliasScript := fmt.Sprintf(`!%s commit`, dinyPath)

	cmd := exec.Command("git", "config", "--global", "alias.auto", aliasScript)
	err = cmd.Run()
	if err != nil {
		ui.Box(ui.BoxOptions{
			Message: fmt.Sprintf("Failed to set git alias: %v", err),
			Variant: ui.Error,
		})
		os.Exit(1)
	}

	ui.Box(ui.BoxOptions{
		Title:   "Git auto alias set up successfully!",
		Message: "Now you can use: git auto",
		Variant: ui.Success,
	})
}

func removeGitAlias() {
	cmd := exec.Command("git", "config", "--global", "--unset", "alias.auto")
	err := cmd.Run()
	if err != nil {
		if strings.Contains(err.Error(), "exit status 5") {
			ui.Box(ui.BoxOptions{
				Message: "No git auto alias found to remove",
				Variant: ui.Warning,
			})
			return
		}
		ui.Box(ui.BoxOptions{
			Message: fmt.Sprintf("Failed to remove git alias: %v", err),
			Variant: ui.Error,
		})
		os.Exit(1)
	}

	ui.Box(ui.BoxOptions{
		Message: "Git alias removed successfully!",
		Variant: ui.Success,
	})
}

func getDinyPath() (string, error) {
	// First try to find diny in PATH
	path, err := exec.LookPath("diny")
	if err == nil {
		return path, nil
	}

	// Try to find diny binary in git root directory
	gitRoot, err := git.FindGitRoot()
	if err == nil {
		dinyPath := filepath.Join(gitRoot, "diny")
		if _, err := os.Stat(dinyPath); err == nil {
			return dinyPath, nil
		}
	}

	// If not in PATH, try to use the current executable
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("diny not found in PATH, git root, and cannot determine executable path: %v", err)
	}

	// Check if this is a temporary path (go run)
	if strings.Contains(execPath, "/tmp/go-build") {
		return "", fmt.Errorf("please build diny first: go build -o diny")
	}

	absPath, err := filepath.Abs(execPath)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(absPath); err != nil {
		return "", err
	}

	return absPath, nil
}

func init() {
	rootCmd.AddCommand(autoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// aliasCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// aliasCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
