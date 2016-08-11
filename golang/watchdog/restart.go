package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// restartPackage will rebuild the given package and restart the process
func restartPackage(packagePath string, arguments []string) {

	fmt.Println("")
	log.Println("BUILDING\n----------------------------")

	// build and install the package
	builder := exec.Command("go", "install", packagePath)

	// redirect all output to the standard console
	builder.Stdout = os.Stdout
	builder.Stderr = os.Stderr

	// run the build command and wait for it to exit
	err := builder.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: could not build the package", err)
		return
	}

	// get the package name to run the package
	executable := getExecutableName(packagePath)

	// stop all previous running instances of the project
	kill := exec.Command("pkill", "-x", executable)

	// redirect all output to the standard console
	kill.Stdout = os.Stdout
	kill.Stderr = os.Stderr

	// note: ignore errors from the kill command
	// (will occur if no executable is running)
	_ = kill.Run()

	// start the go executable
	runner := exec.Command(executable, arguments...)

	// set the current directory to the packagePath
	os.Chdir("/go/src/" + packagePath)

	// redirect all output to the standard console
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr

	// run the program
	err = runner.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: could not run the package", err)
	}

}

// getExecutableName will get the last part of the package path that is used
// to name the go executable
func getExecutableName(packagePath string) string {
	parts := strings.Split(packagePath, "/")
	return parts[len(parts)-1]
}
