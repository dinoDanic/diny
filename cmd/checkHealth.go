package cmd

import (
	"fmt"
	"os"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/ollama"
	"github.com/dinoDanic/diny/ui"
	"github.com/spf13/cobra"
)

var checkHealthCmd = &cobra.Command{
	Use:   "check-health",
	Short: "Check Ollama connection and model availability",
	Long: `Verify that Ollama is running and the configured model is available.

This command checks:
- Connection to the Ollama server
- Availability of the configured model
- Lists all available models on the server

Useful for troubleshooting Ollama configuration issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkOllamaStatus()
	},
}

func checkOllamaStatus() {
	configService := config.GetService()
	if err := configService.LoadUserConfig(); err != nil {
		ui.Box(ui.BoxOptions{
			Message: fmt.Sprintf("Error loading configuration: %v", err),
			Variant: ui.Error,
		})
		os.Exit(1)
	}

	userConfig := configService.GetUserConfig()
	apiConfig := configService.GetAPIConfig()

	if apiConfig.Provider != config.LocalOllama {
		ui.Box(ui.BoxOptions{
			Title:   "Not using Ollama",
			Message: fmt.Sprintf("Current configuration uses: %s\n\nTo use Ollama, run: diny init", apiConfig.Provider),
			Variant: ui.Warning,
		})
		return
	}

	fmt.Println()
	ui.Box(ui.BoxOptions{
		Title:   "Ollama Configuration",
		Message: fmt.Sprintf("URL: %s\nModel: %s", apiConfig.BaseURL, apiConfig.Model),
	})

	fmt.Println()
	ui.Box(ui.BoxOptions{Message: "Checking connection to Ollama..."})

	if err := ollama.CheckHealth(apiConfig.BaseURL); err != nil {
		ui.Box(ui.BoxOptions{
			Title:   "Connection Failed",
			Message: err.Error(),
			Variant: ui.Error,
		})
		os.Exit(1)
	}

	ui.Box(ui.BoxOptions{
		Message: "Connected to Ollama successfully",
		Variant: ui.Success,
	})

	fmt.Println()
	ui.Box(ui.BoxOptions{Message: "Checking if model is available..."})

	if err := ollama.CheckModelExists(apiConfig.BaseURL, apiConfig.Model); err != nil {
		ui.Box(ui.BoxOptions{
			Title:   "Model Check Failed",
			Message: err.Error(),
			Variant: ui.Error,
		})
		os.Exit(1)
	}

	ui.Box(ui.BoxOptions{
		Message: fmt.Sprintf("Model '%s' is available and ready to use", apiConfig.Model),
		Variant: ui.Success,
	})

	fmt.Println()
	ollamaURLSource := config.GetConfigSource("DINY_OLLAMA_URL", userConfig.OllamaURL, "http://127.0.0.1:11434")
	modelSource := config.GetConfigSource("DINY_OLLAMA_MODEL", userConfig.OllamaModel, "llama3.2")

	ui.Box(ui.BoxOptions{
		Title: "Configuration Source",
		Message: fmt.Sprintf("URL source: %s\nModel source: %s", ollamaURLSource, modelSource),
	})
}

func init() {
	rootCmd.AddCommand(checkHealthCmd)
}
