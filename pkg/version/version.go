package version

import (
	"fmt"
	"runtime/debug"
)

func Get() string {
	var revision string
	var modified bool

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}

	if revision == "" {
		return "unavailable"
	}

	if len(revision) > 7 {
		revision = revision[:7]
	}

	if modified {
		return fmt.Sprintf("%s-dirty", revision)
	}

	return revision
}
