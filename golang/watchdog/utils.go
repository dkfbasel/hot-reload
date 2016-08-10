package main

import (
	"strings"
)

// removeLastDirectory will remove the last directory name from the given
// path to allow for direct symlinking
func removeLastDirectory(filepath string) string {

	parts := strings.Split(filepath, "/")
	return strings.Join(parts[:len(parts)-1], "/")

}

// containsAny will check whether any of the matches is part of the given string
func containsAny(source string, matches []string) bool {

	for _, element := range matches {
		if strings.Contains(source, element) == true {
			return true
		}
	}

	return false
}
