package util

import (
	"fmt"
	"strings"
)

var (
	VERSION_NUMBER = fmt.Sprintf("%.02f", 0.01) // the version
	COMMIT         = ""
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
