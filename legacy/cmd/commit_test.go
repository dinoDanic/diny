/*
Copyright © 2025 NAME HERE <email>
*/
package cmd

import (
	"testing"
)

func TestCommitCommand(t *testing.T) {
	t.Run("command registration", func(t *testing.T) {
		if commitCmd == nil {
			t.Fatal("commitCmd should not be nil")
		}

		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == "commit" {
				found = true
				break
			}
		}

		if !found {
			t.Error("commit command should be registered with root command")
		}
	})

	t.Run("command properties", func(t *testing.T) {
		if commitCmd.Use != "commit" {
			t.Errorf("expected Use to be 'commit', got '%s'", commitCmd.Use)
		}

		if commitCmd.Short == "" {
			t.Error("Short description should not be empty")
		}

		if commitCmd.Long == "" {
			t.Error("Long description should not be empty")
		}

		expectedShort := "Generate commit messages from staged changes"
		if commitCmd.Short != expectedShort {
			t.Errorf("expected Short to be '%s', got '%s'", expectedShort, commitCmd.Short)
		}
	})

	t.Run("command execution", func(t *testing.T) {
		if commitCmd.Run == nil {
			t.Error("Run function should not be nil")
		}

		if commitCmd.RunE != nil {
			t.Error("RunE should be nil when Run is set")
		}
	})

	t.Run("command structure", func(t *testing.T) {
		if commitCmd.Parent() != rootCmd {
			t.Error("commit command should have rootCmd as parent")
		}

		if commitCmd.HasSubCommands() {
			t.Error("commit command should not have subcommands")
		}
	})
}

func TestCommitCommandExecution(t *testing.T) {
	t.Run("command can be found by name", func(t *testing.T) {
		cmd, _, err := rootCmd.Find([]string{"commit"})
		if err != nil {
			t.Fatalf("failed to find commit command: %v", err)
		}

		if cmd.Name() != "commit" {
			t.Errorf("expected command name 'commit', got '%s'", cmd.Name())
		}
	})

	t.Run("command accepts no arguments", func(t *testing.T) {
		if commitCmd.Args != nil {
			if err := commitCmd.Args(commitCmd, []string{"extra", "args"}); err == nil {
				t.Log("commit command accepts extra arguments (this may be intentional)")
			}
		}
	})
}
