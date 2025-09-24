package commit

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/helpers"
	"github.com/dinoDanic/diny/ollama"
	"github.com/spf13/cobra"
)

func Main(cmd *cobra.Command, args []string) {
	gitDiffCmd := exec.Command("git", "diff", "--cached",
		"-U0", "--no-color", "--ignore-all-space", "--ignore-blank-lines",
		":(exclude)*.lock", ":(exclude)*package-lock.json", ":(exclude)*yarn.lock",
		":(exclude)node_modules/", ":(exclude)dist/", ":(exclude)build/")

	gitDiff, err := gitDiffCmd.Output()

	if err != nil {
		fmt.Printf("❌ Failed to get git diff: %v\n", err)
		os.Exit(1)
	}

	if len(gitDiff) == 0 {
		fmt.Println("🦴 No staged changes found. Stage files first with `git add`.")
		os.Exit(0)
	}

	cleanDiff := helpers.CleanForAI(string(gitDiff))
	gitDiffLen := len(gitDiff)
	cleanDiffLen := len(cleanDiff)

	userConfig := config.Load()

	systemPrompt := helpers.BuildSystemPrompt(userConfig)
	fullPrompt := systemPrompt + cleanDiff

	fmt.Print("\n")
	if cleanDiffLen > 2000 {
		fmt.Println("⚠️ Large changeset detected — this may take longer to process ⏳")
		fmt.Print("\n")
	}

	if cleanDiffLen == 0 {
		fmt.Println("🌱 No meaningful content detected in the diff.")
		os.Exit(0)
	}
	fmt.Printf("📏 Diff   size → Raw:     %d chars \n", gitDiffLen)
	fmt.Printf("📏 Diff   size → Cleaned: %d chars \n", cleanDiffLen)
	fmt.Printf("📏 Inst   size → Raw:     %d chars \n", len(systemPrompt))
	fmt.Print("\n")

	// Print configuration
	fmt.Println("⚙️  Configuration:")
	fmt.Printf("   • Emoji: %t\n", userConfig.UseEmoji)
	fmt.Printf("   • Conventional: %t\n", userConfig.UseConventional)
	fmt.Printf("   • Tone: %s\n", userConfig.Tone)
	fmt.Printf("   • Length: %s\n", userConfig.Length)
	fmt.Print("\n")
	fmt.Print("🐢 My tiny server is thinking hard, hold tight!")
	fmt.Print("\n")
	fmt.Print("\n")

	commitMessage, err := ollama.MainStream(fullPrompt)

	if err != nil {
		fmt.Printf("💥 Error generating commit message: %v\n", err)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error displaying message: %v\n", err)
	}

	fmt.Print("\n")
	fmt.Print("\n")

	confirmed := confirmPrompt("👉 Do you want to commit with this message?")

	if confirmed {
		commitCmd := exec.Command("git", "commit", "--no-verify", "-m", commitMessage)
		err := commitCmd.Run()
		if err != nil {
			fmt.Printf("❌ Commit failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Commit successfully added to history!")
	} else {
		fmt.Println("🚫 Commit cancelled.")
	}
}

func confirmPrompt(message string) bool {
	var confirmed bool

	err := huh.NewConfirm().
		Title(message).
		Affirmative("Yes").
		Negative("No").
		Value(&confirmed).
		Run()

	if err != nil {
		fmt.Printf("Error running prompt: %v\n", err)
		os.Exit(1)
	}

	return confirmed
}

// ProcessGitDiff processes git diff output and returns cleaned diff and system prompt
// This function is extracted for easier testing
func ProcessGitDiff(gitDiffOutput []byte, userConfig config.UserConfig) (cleanedDiff string, systemPrompt string, err error) {
	if len(gitDiffOutput) == 0 {
		return "", "", fmt.Errorf("no staged changes found")
	}

	cleanDiff := helpers.CleanForAI(string(gitDiffOutput))
	if len(strings.TrimSpace(cleanDiff)) == 0 {
		return "", "", fmt.Errorf("no meaningful content detected in the diff")
	}

	systemPrompt = helpers.BuildSystemPrompt(userConfig)
	return cleanDiff, systemPrompt, nil
}
