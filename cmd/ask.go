/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dinoDanic/diny/ollama"
	"github.com/spf13/cobra"
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Debugg mode",
	Long:  "debug mode",
	Run: func(cmd *cobra.Command, args []string) {
		var prompt string

		// If args provided, use them as prompt
		if len(args) > 0 {
			prompt = strings.Join(args, " ")
		} else {
			// Otherwise, ask for input
			fmt.Print("What would you like to ask? ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("❌ Error reading input: %v\n", err)
				os.Exit(1)
			}
			prompt = strings.TrimSpace(input)
		}

		if prompt == "" {
			fmt.Println("❌ No question provided")
			os.Exit(1)
		}

		fmt.Printf("\n🤔 Thinking about: %s\n\n", prompt)

		response, err := ollama.MainStream(prompt)
		if err != nil {
			fmt.Printf("❌ Error getting response: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n\n✨ Response: %s\n", response)
	},
}

func init() {
	rootCmd.AddCommand(askCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// askCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// askCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
