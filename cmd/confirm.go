package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

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
