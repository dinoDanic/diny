/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/dinoDanic/diny/helpers"
	"github.com/dinoDanic/diny/ollama"
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

		cleanDiff := slimdiff.CleanForAI(string(gitDiff))

		gitDiffLen := len(gitDiff)
		cleanDiffLen := len(cleanDiff)

		fmt.Printf("Raw git diff: %d chars, Clean content: %d chars\n", gitDiffLen, cleanDiffLen)

		if cleanDiffLen > 2000 {
			fmt.Println("Large changeset detected, it will take some time, hold tight!")
		}

		if cleanDiffLen == 0 {
			fmt.Println("No meaningful content changes detected in the diff.")
			os.Exit(1)
		}

		commitMessage, err := ollama.Main(cleanDiff)

		if err != nil {
			fmt.Printf("Error generating commit message: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Generated commit message:\n%s\n", commitMessage)

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
