package commit

import (
	"fmt"
	"os"
	"strings"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/ui"
	"github.com/spf13/cobra"
)

func Main(cmd *cobra.Command, args []string) {
	printMode, _ := cmd.Flags().GetBool("print")

	diff, userConfig := getCommitData(printMode)

	if printMode {
		commitMessage, err := CreateCommitMessage(diff, userConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating commit message: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(strings.TrimSpace(commitMessage))
		return
	}

	configService := config.GetService()
	apiConfig := configService.GetAPIConfig()
	tokens := estimateTokens(diff)

	info := fmt.Sprintf("Diff: ~%d tokens | Config: emoji:%t conv:%t tone:%s len:%s",
		tokens,
		userConfig.UseEmoji,
		userConfig.UseConventional,
		userConfig.Tone,
		userConfig.Length)

	if apiConfig.Provider == config.LocalOllama {
		info += fmt.Sprintf(" | Model: %s", apiConfig.Model)
	}

	ui.Box(ui.BoxOptions{Message: info, Variant: ui.Primary})

	var commitMessage string
	var spinnerMessage string
	if apiConfig.Provider == config.LocalOllama {
		spinnerMessage = "Commit message generating locally..."
	} else {
		spinnerMessage = "Commit message generating via Diny cloud..."
	}

	err := ui.WithSpinner(spinnerMessage, func() error {
		var genErr error
		commitMessage, genErr = CreateCommitMessage(diff, userConfig)
		return genErr
	})

	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("%v", err), Variant: ui.Error})
		os.Exit(1)
	}

	HandleCommitFlow(commitMessage, diff, userConfig)
}

func estimateTokens(text string) int {
	words := len(strings.Fields(text))
	chars := len(text)
	return (words*4 + chars) / 5
}

func getCommitData(isQuietMode bool) (string, *config.UserConfig) {
	gitDiff, err := git.GetGitDiff()

	if err != nil {
		if isQuietMode {
			fmt.Fprintf(os.Stderr, "Failed to get git diff: %v\n", err)
		} else {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to get git diff: %v", err), Variant: ui.Error})
		}
		os.Exit(1)
	}

	if len(gitDiff) == 0 {
		if isQuietMode {
			fmt.Fprintf(os.Stderr, "No staged changes found. Stage files first with `git add`.\n")
		} else {
			ui.Box(ui.BoxOptions{Message: "No staged changes found. Stage files first with `git add`.", Variant: ui.Warning})
		}
		os.Exit(0)
	}

	configService := config.GetService()
	if err := configService.LoadUserConfig(); err != nil {
		if isQuietMode {
			fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		} else {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to load config: %v", err), Variant: ui.Error})
		}
		os.Exit(1)
	}

	userConfig := configService.GetUserConfig()

	return gitDiff, userConfig
}
