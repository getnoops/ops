package version

import (
	"strings"
)

var (
	VERSION_NUMBER = "dev"
	COMMIT         = "unknown"
)

func Version() string {
	parts := []string{VERSION_NUMBER}
	if COMMIT != "" {
		parts = append(parts, COMMIT)
	}

	return strings.Join(parts, " ")
}

func Commit() string {
	return COMMIT
}
