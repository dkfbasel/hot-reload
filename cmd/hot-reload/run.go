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

// runBuild will run a build command and restart the package
func runBuild(config Config) {

	// do only allow one build command at the time
	mtx.Lock()
	defer mtx.Unlock()

	// kill any current runner
	if runner != nil {
		runner.Process.Kill() // nolint:errcheck
	}

	log.Println("BUILDING\n----------------------------")

	builder := exec.Command("go", "build", "-o", "/tmp/app")
	builder.Dir = config.Directory
	builder.Stdout = os.Stdout
	builder.Stderr = os.Stderr
	err := builder.Run()
	if err != nil {
		log.Printf("[ERROR] %+v", err)
		return
	}

	log.Println("RUNNING\n----------------------------")

	// executing the binary
	runner = exec.Command("/tmp/app", config.Arguments...)
	runner.Dir = config.Directory
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	err = runner.Start()
	if err != nil {
		log.Printf("[ERROR] %+v", err)
	}

	<-time.After(time.Millisecond * 300)

}

// runTest will run go test on the package
func runTest(config Config) {

	// do only allow one build command at the time
	mtx.Lock()
	defer mtx.Unlock()

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
		fmt.Fprintln(os.Stderr, "ERROR: could not run go test", err)
	}

	return

}
