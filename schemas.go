package revisorschemas

import (
	"embed"
	"runtime/debug"
)

var version = "v0.0.0-dev"

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	for _, m := range info.Deps {
		if m.Path != "github.com/ttab/revisorschemas" {
			continue
		}

		version = m.Version
		break
	}
}

//go:embed *.json
var specifications embed.FS

func Version() string {
	return version
}

func Files() embed.FS {
	return specifications
}
