/* Copyright © 2025 dinoDanic dino.danic@gmail.com */
package version

import (
	"runtime/debug"
)

func Get() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}
	return "dev"
}
