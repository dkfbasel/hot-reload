package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const webpackConfigFile = "webpack.config.js"

func main() {

	// define the initial directory that we are going to search through
	var directory string
	flag.StringVar(&directory, "directory", "", "the path of the webpack directory")

	if directory == "" {
		directory = os.Getenv("DIRECTORY")
	}

	if directory == "" {
		directory = findWebpackDirectory("../../_test")
	} else {
		directory = "/app/" + directory
	}

	if directory == "" {
		fmt.Println("No webpack directory found. Please specify the directory using -e \"DIRECTORY=web\"")
	}

	err := symlinkGlobalNodeModules(directory)
	if err != nil {
		fmt.Println("could not symlink the global node modules:", err)
		return
	}

	err = startWebpack(directory)
	if err != nil {
		fmt.Println("could not start the webpack dev server:", err)
		return
	}

}

// findWebpackDirectory will try to find the webpack directory by looking for
// a webpack.config.js file.
func findWebpackDirectory(searchDirectory string) string {

	var webpackDirectory string

	// find the directory that contains a webpack config file
	err := filepath.Walk(searchDirectory, func(filePath string, f os.FileInfo, err error) error {

		// try to find the webpack directory file
		if strings.Contains(filePath, webpackConfigFile) {
			webpackDirectory = path.Dir(filePath)
			return filepath.SkipDir
		}

		excludeDirs := []string{"/vendor", "/node_modules", ".git", ".svn"}
		if containsAny(filePath, excludeDirs) {
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}

	return webpackDirectory
}

// containsAny will check whether any of the matches is part of the given
// string.
func containsAny(source string, matches []string) bool {

	for _, element := range matches {
		if strings.Contains(source, element) == true {
			return true
		}
	}

	return false
}

// symlinkGlobalNodeModules will link the globally installed node modules in the
// container into the node_modules folder of the project.
func symlinkGlobalNodeModules(directory string) error {

	// define the path the node_modules directory
	nodeModules := directory + "/node_modules"

	// check if a node directory exists in the same directory
	// or create if necessary
	if _, err := os.Stat(nodeModules); os.IsNotExist(err) {
		os.Mkdir(nodeModules, 0777)
	}

	// symlink the global node modules into the directory
	symlink := exec.Command("ln", "-s", "-f", "/usr/local/lib/node_modules/*", nodeModules)

	// redirect all output to the standard console
	symlink.Stdout = os.Stdout
	symlink.Stderr = os.Stderr

	err := symlink.Run()
	return err

}

// startWebpack will try to start a webpack dev server using the command
// npm run dev.
func startWebpack(directory string) error {
	// start the webpack dev server using an npm run command
	webpack := exec.Command("npm", "run", "dev")

	// set the current directory to the webpack directory
	os.Chdir(directory)

	// redirect all output to the standard console
	webpack.Stdout = os.Stdout
	webpack.Stderr = os.Stderr

	// try to run the webpack service
	err := webpack.Run()
	return err

}
