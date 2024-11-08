package postman

import (
	"fmt"
	"regexp"
)

const separator = "/"

var fre = regexp.MustCompile(fmt.Sprintf(`%s+([^%s]+)`, separator, separator))

func SlicePath(path string) []string {
	matches := fre.FindAllStringSubmatch(path, -1)

	var paths []string
	for _, m := range matches {
		paths = append(paths, m[1])
	}

	return paths
}
