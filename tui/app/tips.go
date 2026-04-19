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
	"press p to commit and push in a single step",
	"press n to commit with --no-verify and skip git hooks",
	"press r to regenerate the message from scratch",
	"press v to cycle through alternate message variants",
	"press f to give feedback and refine the current message",
	"run `diny timeline` to summarize recent commits for a standup or PR body",
	"run `diny changelog` to draft release notes from git history",
	"run `diny link lazygit` to wire diny into your lazygit keybindings",
	"run `diny config` to tweak tone, length, emoji, and custom instructions",
	"set custom_instructions in `diny config` to steer tone across every message",
}

func randomTip() string {
	return tips[rand.Intn(len(tips))]
}
