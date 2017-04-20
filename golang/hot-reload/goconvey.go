package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// runGoconvey will start goconvey in the main directory
func runGoconvey(packagePath string, arguments []string) {

	fmt.Println("")
	log.Println("START GOCONVEY\n----------------------------")

	// set the current directory to the packagePath
	os.Chdir("/go/src/" + packagePath) // nolint: errcheck

	// append the host to all other arguments (must be 0.0.0.0 to be exposed
	// by docker
	arguments = append([]string{"-host", "0.0.0.0"}, arguments...)

	// start the go executable
	goconvey := exec.Command("goconvey", arguments...)

	// redirect all output to the standard console
	goconvey.Stdout = os.Stdout
	goconvey.Stderr = os.Stderr

	// run the program
	err := goconvey.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: could not start goconvey", err)
	}

}
