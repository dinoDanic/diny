/*
Copyright © 2025 NAME HERE <email>
*/

package commit

import (
	"testing"

	"github.com/dinoDanic/diny/config"
)

// TestMain is commented out because Main() calls os.Exit() which causes test failures
// In a production codebase, you'd want to refactor Main to return errors instead of calling os.Exit()
// For now, we test the extracted ProcessGitDiff function which contains the core logic
//
// func TestMain(t *testing.T) {
// 	// This test is skipped because Main() calls os.Exit()
// 	// which terminates the test process
// }

func TestProcessGitDiff(t *testing.T) {
	userConfig := config.UserConfig{
		UseEmoji:        true,
		UseConventional: false,
		Tone:            config.Casual,
		Length:          config.Normal,
	}

	t.Run("Empty diff should return error", func(t *testing.T) {
		_, _, err := ProcessGitDiff([]byte{}, userConfig)
		if err == nil {
			t.Error("Expected error for empty diff, got nil")
		}
	})

	t.Run("Valid diff should process successfully", func(t *testing.T) {
		sampleDiff := `diff --git a/test.go b/test.go
index 1234567..abcdefg 100644
--- a/test.go
+++ b/test.go
@@ -1,3 +1,4 @@
 func main() {
+    fmt.Println("hello")
     return
 }`

		cleanDiff, systemPrompt, err := ProcessGitDiff([]byte(sampleDiff), userConfig)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if cleanDiff == "" {
			t.Error("Expected non-empty cleaned diff")
		}

		if systemPrompt == "" {
			t.Error("Expected non-empty system prompt")
		}

		t.Logf("Cleaned diff: %s", cleanDiff)
		t.Logf("System prompt: %s", systemPrompt)
	})
}
