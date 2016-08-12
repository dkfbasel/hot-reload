package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

// createSymlinkForPackage will create a new symlink from the app directory to
// the specified package name in the go source directory
func createSymlinkForPackage(gopackagePath string) {

	// create the directories in the path to our package. note that we remove the
	// last directory to allow symlinking to work as expected
	mkdir := exec.Command("mkdir", "-p", "/go/src/"+removeLastDirectory(gopackagePath))

	// redirect all output to the standard console
	mkdir.Stdout = os.Stdout
	mkdir.Stderr = os.Stderr

	err := mkdir.Run()
	if err != nil {
		log.Fatalf("mkdir command finished with error: %s\n", err)
	}

	// link our directory into the go src directory
	symlink := exec.Command("ln", "-s", "-f", "/app/", "/go/src/"+gopackagePath)

	// redirect all output to the standard console
	symlink.Stdout = os.Stdout
	symlink.Stderr = os.Stderr

	err = symlink.Run()
	if err != nil {
		log.Fatalf("ln command finished with error: %s\n", err)
	}

}

// removeLastDirectory will remove the last directory name from the given
// path to allow for direct symlinking
func removeLastDirectory(filepath string) string {

	parts := strings.Split(filepath, "/")
	return strings.Join(parts[:len(parts)-1], "/")

}
