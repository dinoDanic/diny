/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"diny/ollama"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate AI-powered commit messages from staged changes",
	Long: `Generate intelligent commit messages using AI based on your staged git changes.
	
This command analyzes your staged git changes and uses AI to generate a concise,
well-formatted commit message following best practices.

Example:
  diny commit`,
	Run: func(cmd *cobra.Command, args []string) {
		// Use optimized git diff with built-in filtering
		gitDiffCmd := exec.Command("git", "diff", "--cached",
			"-U0", "--no-color", "--ignore-all-space", "--ignore-blank-lines",
			"--diff-filter=AM", // Only Added and Modified files
			":(exclude)*.lock", ":(exclude)*package-lock.json", ":(exclude)*yarn.lock",
			":(exclude)*.min.js", ":(exclude)*.min.css", ":(exclude)*.bundle.js",
			":(exclude)*.jpg", ":(exclude)*.jpeg", ":(exclude)*.png", ":(exclude)*.gif",
			":(exclude)*.pdf", ":(exclude)*.zip", ":(exclude)*.exe", ":(exclude)*.dll",
			":(exclude)node_modules/", ":(exclude)dist/", ":(exclude)build/")

		gitDiff, err := gitDiffCmd.Output()

		if err != nil {
			fmt.Printf("Error getting git diff: %v\n", err)
			os.Exit(1)
		}

		if len(gitDiff) == 0 {
			fmt.Println("No staged changes found. Please stage your changes with 'git add' first.")
			os.Exit(1)
		}

		gitDiffLen := len(gitDiff)

		fmt.Print("git diff len: ", gitDiffLen)
		fmt.Print("\n")
		if gitDiffLen > 5000 {
			fmt.Print("ohh this will take some time.. will optimaze it! hold tight!")
		}
		fmt.Print("\n")

		// NOTE: this is lame, should optimaze somehow
		// FALLBACK TO NAME ONLY
		// Fallback to file names for very large diffs
		// if len(gitDiff) > 8000 {
		// 	stagedFilesCmd := exec.Command("git", "diff", "--cached", "--name-only")
		// 	stagedFiles, err := stagedFilesCmd.Output()
		// 	if err == nil && len(stagedFiles) > 0 {
		// 		fmt.Println("Large changeset detected, using file summary.")
		// 		files := strings.TrimSpace(string(stagedFiles))
		// 		commitMessage, err := ollama.Main(fmt.Sprintf("Files modified: %s", files))
		// 		if err != nil {
		// 			fmt.Printf("Error generating commit message: %v\n", err)
		// 			os.Exit(1)
		// 		}
		// 		fmt.Printf("Generated commit message:\n%s\n", commitMessage)
		// 		return
		// 	}
		// }

		commitMessage, err := ollama.Main(string(gitDiff))
		if err != nil {
			fmt.Printf("Error generating commit message: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Generated commit message:\n%s\n", commitMessage)

		// Ask user for confirmation using Bubbletea
		confirmed := confirmPrompt("Do you want to commit with this message?")

		if confirmed {
			commitCmd := exec.Command("git", "commit", "-m", commitMessage)
			err := commitCmd.Run()
			if err != nil {
				fmt.Printf("Error committing: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Committed successfully!")
		} else {
			fmt.Println("❌ Commit cancelled.")
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
