package util

import (
	"fmt"
	"strings"
)

var (
	VERSION_NUMBER = fmt.Sprintf("%.02f", 0.01)
	COMMIT         = ""
)

func Version() string {
	parts := []string{VERSION_NUMBER}
	if COMMIT != "" {
		parts = append(parts, COMMIT)
	}

	return strings.Join(parts, " ")
}
