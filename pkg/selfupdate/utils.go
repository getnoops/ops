package selfupdate

import (
	"fmt"
	"strings"
)

func IsDifferent(commit string, filename string) (bool, error) {
	splits := strings.Split(filename, "-")
	if len(splits) != 4 {
		return false, fmt.Errorf("invalid filename: %s", filename)
	}

	short := commit
	if len(commit) > 8 {
		short = commit[:8]
	}

	if splits[1] != short {
		return true, nil
	}

	return false, nil
}
