package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/ollama"
	"github.com/dinoDanic/diny/ui"
	"github.com/spf13/cobra"
)

func promptConfigScope() bool {
	var scopeChoice string

	err := huh.NewSelect[string]().
		Title("Configuration scope").
		Description("Choose where to save your preferences").
		Options(
			huh.NewOption("Global - Apply to all repositories", "global"),
			huh.NewOption("Local - Only for this repository", "local"),
		).
		Value(&scopeChoice).
		Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	return scopeChoice == "global"
}

func promptConfigAction(isGlobal bool, configExists bool) string {
	var action string
	scopeName := "local"
	if isGlobal {
		scopeName = "global"
	}

	var title string
	var options []huh.Option[string]

	if configExists {
		title = fmt.Sprintf("%s configuration already exists", scopeName)
		options = []huh.Option[string]{
			huh.NewOption("Update/Edit - Modify existing values", "update"),
			huh.NewOption("Create New - Start fresh", "new"),
		}
	} else {
		title = fmt.Sprintf("Create %s configuration", scopeName)
		options = []huh.Option[string]{
			huh.NewOption("Update/Edit - Use existing as base", "update"),
			huh.NewOption("Create New - Start fresh", "new"),
		}
	}

	err := huh.NewSelect[string]().
		Title(title).
		Description("What would you like to do?").
		Options(options...).
		Value(&action).
		Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	return action
}

func RunConfigurationSetup(existingConfig *config.UserConfig) config.UserConfig {
	var userConfig config.UserConfig

	if existingConfig != nil {
		userConfig = *existingConfig
	} else {
		userConfig = config.UserConfig{
			UseEmoji:        false,
			UseConventional: false,
			UseLocalAPI:     false,
			Tone:            config.Casual,
			Length:          config.Short,
		}
	}

	err := huh.NewConfirm().
		Title("Use emoji prefixes in commit messages?").
		Description("Add emojis like ✨ feat: or 🐛 fix: to commit messages").
		Affirmative("Yes").
		Negative("No").
		Value(&userConfig.UseEmoji).
		Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = huh.NewConfirm().
		Title("Use Conventional Commits format?").
		Description("Format: type(scope): description").
		Affirmative("Yes").
		Negative("No").
		Value(&userConfig.UseConventional).
		Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	var apiChoice string
	if userConfig.UseLocalAPI {
		apiChoice = "custom-ollama"
	} else {
		apiChoice = "diny-cloud"
	}

	err = huh.NewSelect[string]().
		Title("Choose AI backend").
		Description("Select where commit messages will be generated").
		Options(
			huh.NewOption("Diny Cloud - Default (https://diny-cli.vercel.app)", "diny-cloud"),
			huh.NewOption("Custom Ollama - Configure your own instance", "custom-ollama"),
		).
		Value(&apiChoice).
		Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	switch apiChoice {
	case "diny-cloud":
		userConfig.UseLocalAPI = false

	case "custom-ollama":
		userConfig.UseLocalAPI = true

		ollamaURL := "http://127.0.0.1:11434"
		if existingConfig != nil && existingConfig.OllamaURL != "" {
			ollamaURL = existingConfig.OllamaURL
		}

		err = huh.NewInput().
			Title("Ollama URL").
			Description("Edit or press Enter to keep default").
			Value(&ollamaURL).
			Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		userConfig.OllamaURL = ollamaURL

		fmt.Println()
		ui.Box(ui.BoxOptions{Message: "Checking connection to Ollama..."})
		if err := validateOllamaConnection(ollamaURL); err != nil {
			ui.Box(ui.BoxOptions{Message: err.Error(), Variant: ui.Warning})

			var continueAnyway bool
			err = huh.NewConfirm().
				Title("Continue anyway?").
				Description("Connection failed. Save config and fix later?").
				Affirmative("Yes, save anyway").
				Negative("No, go back").
				Value(&continueAnyway).
				Run()

			if err != nil || !continueAnyway {
				ui.Box(ui.BoxOptions{Message: "Configuration cancelled", Variant: ui.Error})
				os.Exit(1)
			}
		} else {
			ui.Box(ui.BoxOptions{Message: "Connected to Ollama successfully", Variant: ui.Success})
		}

		ollamaModel := "llama3.2"
		if existingConfig != nil && existingConfig.OllamaModel != "" {
			ollamaModel = existingConfig.OllamaModel
		}

		err = huh.NewInput().
			Title("Ollama Model").
			Description("Edit or press Enter to keep default").
			Value(&ollamaModel).
			Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		userConfig.OllamaModel = ollamaModel

		fmt.Println()
		ui.Box(ui.BoxOptions{Message: "Checking if model exists..."})
		if err := validateOllamaModel(ollamaURL, ollamaModel); err != nil {
			ui.Box(ui.BoxOptions{Message: err.Error(), Variant: ui.Warning})

			var continueAnyway bool
			err = huh.NewConfirm().
				Title("Continue anyway?").
				Description("Model check failed. Save config and fix later?").
				Affirmative("Yes, save anyway").
				Negative("No, go back").
				Value(&continueAnyway).
				Run()

			if err != nil || !continueAnyway {
				ui.Box(ui.BoxOptions{Message: "Configuration cancelled", Variant: ui.Error})
				os.Exit(1)
			}
		} else {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Model '%s' is available", ollamaModel), Variant: ui.Success})
		}
	}

	err = huh.NewSelect[config.Tone]().
		Title("Choose your commit message tone").
		Options(
			huh.NewOption("Professional - formal and matter-of-fact", config.Professional),
			huh.NewOption("Casual - light but clear", config.Casual),
			huh.NewOption("Friendly - warm and approachable", config.Friendly),
		).
		Value(&userConfig.Tone).
		Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = huh.NewSelect[config.Length]().
		Title("Choose your commit message length").
		Options(
			huh.NewOption("Short - subject only (no body)", config.Short),
			huh.NewOption("Normal - subject + optional body (1-4 bullets)", config.Normal),
			huh.NewOption("Long - subject + detailed body (2-6 bullets)", config.Long),
		).
		Value(&userConfig.Length).
		Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	return userConfig
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Diny configuration with an interactive setup",
	Long: `Initialize Diny configuration with an interactive setup.

This command will guide you through configuring your commit message preferences:
- Scope: Configure globally for all repositories or locally for this repository
- Emoji: Add emoji prefixes to commit messages
- Format: Conventional commits or free-form messages
- Local API: Use local Ollama instance (requires installation)
- Tone: Professional, casual, or friendly
- Length: Short, normal, or detailed messages

Config precedence: local > global > defaults`,
	Run: func(cmd *cobra.Command, args []string) {
		globalFlag, _ := cmd.Flags().GetBool("global")

		var isGlobal bool

		if cmd.Flags().Changed("global") {
			isGlobal = globalFlag
		} else {
			isGlobal = promptConfigScope()
		}

		var existingConfig *config.UserConfig
		var configExists bool
		var err error

		if isGlobal {
			existingConfig, err = config.LoadGlobal()
			configExists = existingConfig != nil && err == nil
		} else {
			existingConfig, err = config.LoadLocal()
			configExists = existingConfig != nil && err == nil
		}

		action := promptConfigAction(isGlobal, configExists)

		switch action {
		case "update":
			fmt.Println()
			ui.RenderTitle("Update/Edit Configuration")
			if configExists {
				fmt.Println("Current values are pre-filled. Press Enter to keep, or change them.")
			} else {
				if isGlobal {
					existingConfig, _ = config.LoadLocal()
				} else {
					existingConfig, _ = config.LoadGlobal()
				}
				if existingConfig != nil {
					fmt.Println("Using existing configuration as base. Modify as needed.")
				} else {
					fmt.Println("No existing configuration found. Using defaults.")
				}
			}
			fmt.Println()
		case "new":
			existingConfig = nil
			fmt.Println()
			ui.RenderTitle("Create New Configuration")
			fmt.Println()
		}

		userConfig := RunConfigurationSetup(existingConfig)

		if isGlobal {
			err = config.SaveGlobal(userConfig)
			if err != nil {
				fmt.Printf("Error saving global configuration: %v\n", err)
				os.Exit(1)
			}

			globalPath, _ := config.GetGlobalConfigPath()
			fmt.Println()
			config.PrintConfiguration(userConfig)
			ui.Box(ui.BoxOptions{
				Title:   "Global configuration saved!",
				Message: fmt.Sprintf("Location: %s\n\nAll repositories will use these settings by default", globalPath),
				Variant: ui.Success,
			})
		} else {
			err = config.SaveLocal(userConfig)
			if err != nil {
				fmt.Printf("Error saving local configuration: %v\n", err)
				os.Exit(1)
			}

			localPath, _ := config.GetLocalConfigPath()
			fmt.Println()
			config.PrintConfiguration(userConfig)
			ui.Box(ui.BoxOptions{
				Title:   "Local configuration saved!",
				Message: fmt.Sprintf("Location: %s\n\nThis repository will use these settings", localPath),
				Variant: ui.Success,
			})
		}
	},
}

func validateOllamaConnection(baseURL string) error {
	return ollama.CheckHealth(baseURL)
}

func validateOllamaModel(baseURL, modelName string) error {
	return ollama.CheckModelExists(baseURL, modelName)
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP("global", "g", false, "Skip prompt and save configuration globally for all repositories")
}
