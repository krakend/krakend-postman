package postman

import (
	"strings"
)

const separator = "/"

func slicePath(path string) []string {
	var paths []string
	for _, p := range strings.Split(path, separator) {
		clean := strings.Trim(p, " ")
		if len(clean) > 0 {
			paths = append(paths, p)
		}
	}
	return paths
}
