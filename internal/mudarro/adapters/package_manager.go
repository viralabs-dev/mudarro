package adapters

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Exact SemVer 2.0 syntax; build metadata (including Corepack hash text) is
// accepted syntactically, not verified or resolved by scanning.
var exactManagerVersion = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`)

func packageManagerName(value string) (string, error) {
	name, version, versioned := strings.Cut(value, "@")
	known := name == "npm" || name == "pnpm" || name == "yarn" || name == "bun"
	valid := known && strings.IndexFunc(value, unicode.IsSpace) < 0
	if versioned {
		valid = valid && exactManagerVersion.MatchString(version)
		core, _, _ := strings.Cut(version, "+")
		_, prerelease, hasPrerelease := strings.Cut(core, "-")
		if hasPrerelease {
			for _, part := range strings.Split(prerelease, ".") {
				numeric := part != "" && strings.IndexFunc(part, func(r rune) bool { return r < '0' || r > '9' }) < 0
				if numeric && len(part) > 1 && part[0] == '0' {
					valid = false
				}
			}
		}
	}
	if !valid {
		return "", fmt.Errorf("invalid packageManager %q: use npm, pnpm, yarn or bun, optionally with an exact SemVer version; custom tools require explicit commands", value)
	}
	return name, nil
}
