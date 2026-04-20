package version

// Version is set at build time via ldflags
// (-X github.com/dinoDanic/diny/version.Version=...).
var Version = "dev"

func Get() string {
	return Version
}
