package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// createSymlinkForPackage will create a new symlink from the app directory to
// the specified package name in the go source directory
func createSymlinkForPackage(config Config) {

	// define the target directory to put the source code into (in the go src directory)
	targetDirectory := "/go/src/" + config.ProjectPath

	// create the directories in the path to our package
	mkdir := exec.Command("mkdir", "-p", targetDirectory)

	// redirect all output to the standard console
	mkdir.Stdout = os.Stdout
	mkdir.Stderr = os.Stderr

	err := mkdir.Run()
	if err != nil {
		log.Fatalf("mkdir command finished with error: %s\n", err)
	}

	// define the source directory for the symlinking
	sourceDirectory := "/app"

	// get all directories in the global node module directory
	content, err := ioutil.ReadDir(sourceDirectory)
	if err != nil {
		fmt.Println("Error reading the application directory:\n", err)
	}

	for _, item := range content {
		// symlink the global node modules into the directory
		symlink := exec.Command("ln", "-s", "-f", sourceDirectory+"/"+item.Name(), targetDirectory)

		// redirect all output to the standard console
		symlink.Stdout = os.Stdout
		symlink.Stderr = os.Stderr

		err := symlink.Run()
		if err != nil {
			fmt.Printf("Could not symlink %s.\n%s", item.Name(), err)
		}
	}

}
