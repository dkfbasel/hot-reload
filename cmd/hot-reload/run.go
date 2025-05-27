package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

var mtx sync.Mutex = sync.Mutex{}
var runner *exec.Cmd
var builder *exec.Cmd

// runBuild will run a build command and restart the package
func runBuild(config Config) {

	// kill any current builder
	if builder != nil {
		builder.Process.Kill() // nolint:errcheck
		builder = nil
	}

	mtx.Lock()
	defer mtx.Unlock()

	fmt.Printf("\n%s RELOAD\n----------------------------\n", time.Now().Format("2006/01/02 15:04:05"))

	builder = exec.Command("go", "build", "-o", "/tmp/app")
	builder.Dir = config.Directory
	builder.Stdout = os.Stdout
	builder.Stderr = os.Stderr
	err := builder.Run()
	if err != nil {
		log.Printf("[BUILD ERROR] %+v", err)
		return
	}

	// kill any current runner
	if runner != nil {
		runner.Process.Kill() // nolint:errcheck
		runner = nil
	}

	// executing the binary using delve debugger
	args := []string{
		"exec", "/tmp/app",
		"--headless", "--listen=:2345", "--api-version=2", "--accept-multiclient",
	}

	// append additional arguments to the debugged application
	if len(config.Arguments) > 0 {
		args = append(args, "--")
		args = append(args, config.Arguments...)
	}

	runner = exec.Command("dlv", args...)
	runner.Dir = config.Directory
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	err = runner.Start()
	if err != nil {
		log.Printf("[RUN ERROR] %+v", err)
	}

}

// runTest will run go test on the package
func runTest(config Config) {

	// kill any current runner
	if runner != nil {
		runner.Process.Kill() // nolint:errcheck
	}

	// clear the screen before each new run
	log.Println("TESTING\n----------------------------")

	// testing
	arguments := append([]string{"test"}, config.Arguments...)
	runner = exec.Command("go", arguments...)
	runner.Dir = config.Directory
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr

	// run the program
	err := runner.Run()
	if err != nil {
		log.Printf("[TEST ERROR] %+v", err)
	}

}
