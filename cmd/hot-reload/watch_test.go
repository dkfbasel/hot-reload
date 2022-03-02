package main

import (
	"fmt"
	"testing"
)

func TestMatchesAny(t *testing.T) {

	rootDir := "/app"

	tests := []Test{
		{
			// match directory in app directly
			path:     "/app/dir",
			patterns: []string{"dir"},
			match:    true,
		},
		{
			// do not match directory in subdirectory
			path:     "/app/abc/dir",
			patterns: []string{"dir"},
			match:    false,
		},
		{
			// match directory with wildcard matching
			path:     "/app/abc/dir",
			patterns: []string{"*dir"},
			match:    true,
		},
		{
			// match directory with wildcard matching
			path:     "/app/abc/xyz/dir",
			patterns: []string{"*dir"},
			match:    true,
		},
		{
			// match file in root directory
			path:     "/app/file.log",
			patterns: []string{"file.log"},
			match:    true,
		},
		{
			// match file with wildcard in root directory
			path:     "/app/file.log",
			patterns: []string{"*.log"},
			match:    true,
		},
		{
			// match file with wildcard in subdirectory
			path:     "/app/abc/file.log",
			patterns: []string{"*.log"},
			match:    true,
		},
		{
			// match file with wildcard in subsubdirectory
			path:     "/app/abc/xzy/file.log",
			patterns: []string{"*.log"},
			match:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {

			for i := range test.patterns {
				test.patterns[i] = fmt.Sprintf("%s/%s", rootDir, test.patterns[i])
			}

			result := matchesAny(test.path, test.patterns)
			if test.match != result {
				t.Fail()
			}
		})
	}

}

type Test struct {
	path     string
	patterns []string
	match    bool
}
