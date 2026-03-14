package app

import "math/rand"

var tips = []string{
	"use [ and ] to browse previously generated messages in this session",
	"press t to force a conventional commit type (feat, fix, docs...)",
	"press L to cycle message length: short → normal → long",
	"press A to regenerate from HEAD diff and amend on commit",
	"press M to toggle emoji on/off for this session",
	"press s to save as a draft — useful with lazygit",
	"press d to view the full staged diff before committing",
	"press y to copy the commit message to your clipboard",
	"press e to edit the message inline without leaving diny",
	"press E to edit the message inside your default editor",
	"press ? to see all available keyboard shortcuts",
}

func randomTip() string {
	return tips[rand.Intn(len(tips))]
}
