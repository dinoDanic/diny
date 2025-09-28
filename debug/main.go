package main

import (
	"github.com/dinoDanic/diny/ui"
	"github.com/dinoDanic/diny/update"
)

func main() {
	ui.DebugUI()

	checker := update.NewUpdateChecker("v0.1.0")
	checker.CheckForUpdate()
}
