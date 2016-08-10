package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	// parse the gopackage from the command line or an environment variable.
	// note: a tmp directory will be used if no project path is supplied, this
	// might break some import statements that point to the project directory
	projectPath, subdirectory := parseProjectInfo()

	// print the package path for the user
	fmt.Printf("PACKAGE:\t%s\n", projectPath)
	fmt.Printf("DIRECTORY:\t%s\n", subdirectory)

	if projectPath == "tmp" {
		fmt.Printf("please note, that import paths in the project directory will probably not work as intended")
	}

	// create a symlink for the package to allow for compilation
	createSymlinkForPackage(projectPath)

	// rebuild and start the package
	restartPackage(projectPath + subdirectory)

	// watch the package directory for changes
	watchPackageDirectory(projectPath + subdirectory)
}

// parseProjectInfo will parse the necessary external information from the command line
// or the environment and return an error if the flag is not defined
func parseProjectInfo() (string, string) {

	// the package import path should be supplied via flag or environment variable
	var projectPath string
	var subdirectory string
	flag.StringVar(&projectPath, "project", "", "the path of the project relative to the go source directory")
	flag.StringVar(&subdirectory, "directory", "", "(optional) relative path to the directory of the current go package to watch")
	flag.Parse()

	// try to parse the data from the environment if not supplied via flag
	if projectPath == "" {
		projectPath = os.Getenv("PROJECT")
	}

	if subdirectory == "" {
		subdirectory = os.Getenv("DIRECTORY")
	}

	// set the project path to tmp if not specified -> this will break package
	// imports that are contained in the same project directory
	if projectPath == "" {
		projectPath = "development/tmp"
	}

	// ensure that the subdirectory starts with a slash
	if subdirectory != "" && strings.HasPrefix(subdirectory, "/") == false {
		subdirectory = "/" + subdirectory
	}

	return projectPath, subdirectory

}
