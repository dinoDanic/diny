package update

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/ui"
)

// PromptIfAvailable consumes the channel returned by CheckAsync with a short
// timeout and, if a newer version is available, asks the user whether to
// update. A "no" snoozes the prompt for 48h.
func (uc *UpdateChecker) PromptIfAvailable(ch <-chan string) {
	var latest string
	select {
	case latest = <-ch:
	case <-time.After(300 * time.Millisecond):
		return
	}

	if latest == "" {
		return
	}
	if !uc.compareVersions(uc.currentVersion, latest) {
		return
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
		return
	}

	if !confirmed {
		_ = dismiss()
		return
	}

	if err := uc.PerformUpdate(); err != nil {
		ui.Error("Update failed: %v", err)
	}
}
