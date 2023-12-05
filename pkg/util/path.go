package util

import (
	"os/user"
	"path/filepath"
	"strings"
)

func ResolvePath(path string) (string, error) {
	p := path

	if strings.HasPrefix(p, "~") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		p = dir + p[1:]
	}

	return filepath.Abs(p)
}
