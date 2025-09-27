package commit

import (
	"github.com/dinoDanic/diny/ui"
)

func RenderTitle(text string) string {
	return ui.RenderTitle(text)
}

func RenderSuccess(text string) string {
	return ui.RenderSuccess(text)
}

func RenderError(text string) string {
	return ui.RenderError(text)
}

func RenderWarning(text string) string {
	return ui.RenderWarning(text)
}

func RenderInfo(text string) string {
	return ui.RenderInfo(text)
}

func RenderCommitMessage(message string) string {
	return ui.RenderBox("Commit Message", message)
}

func RenderNote(note string) string {
	return ui.RenderNote(note)
}

func RenderStep(step string) string {
	return ui.RenderStep(step)
}

func WithSpinner(message string, fn func() error) error {
	return ui.WithSpinner(message, fn)
}
