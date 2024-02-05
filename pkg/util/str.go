package util

import "strings"

func JoinStrings[T any](elems []T, selector func(T) string, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return selector(elems[0])
	}

	var b strings.Builder
	b.WriteString(selector(elems[0]))
	for _, elem := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(selector(elem))
	}
	return b.String()
}
