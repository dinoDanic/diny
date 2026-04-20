package update

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/ui"
)

// PromptIfAvailable consumes the channel returned by CheckAsync with a short
// timeout and, if a newer version is available, asks the user whether to
// update. Returns true when the user was shown the update prompt.
func (uc *UpdateChecker) PromptIfAvailable(ch <-chan string) bool {
	var latest string
	select {
	case latest = <-ch:
	case <-time.After(300 * time.Millisecond):
		return false
	}

	if latest == "" {
		return false
	}
	if !uc.compareVersions(uc.currentVersion, latest) {
		return false
	}

	fmt.Println()

	var confirmed bool
	err := huh.NewConfirm().
		Title(fmt.Sprintf("Update diny to %s?", latest)).
		Affirmative("Yes").
		Negative("No").
		Value(&confirmed).
		Run()
	if err != nil {
		return true
	}

	if !confirmed {
		return true
	}

	if err := uc.PerformUpdate(); err != nil {
		ui.Error("Update failed: %v", err)
	}
	return true
}
