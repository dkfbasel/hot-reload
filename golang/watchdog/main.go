package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/exp/inotify"
)

func main() {

	// get the package name from the environment or the command-line
	var gopackage string
	flag.StringVar(&gopackage, "gopackage", "", "import path of the go package to watch")
	flag.Parse()

	// check if a package has been supplied
	if gopackage == "" {
		gopackage = os.Getenv("GOPACKAGE")

		if gopackage == "" || gopackage == "undefined" {
			fmt.Println("ERROR: A go package import path is required.\n-e \"GOPACKAGE=bitbucket.com/dkfbasel/test\"\n\n")
			return
		}
	}

	// Print package name for the user
	fmt.Printf("PACKAGE: %s\n\n", gopackage)

	// remove the last directory to allow direct symlinking
	withoutLastDirectory := removeLastDirectory(gopackage)

	// create the directories in the path to our package
	mkdir := exec.Command("mkdir", "-p", "/go/src/"+withoutLastDirectory)

	// redirect all output to the standard console
	mkdir.Stdout = os.Stdout
	mkdir.Stderr = os.Stderr

	err := mkdir.Run()
	if err != nil {
		log.Fatalf("Command finished with error: %s\n", err)
	}

	// link our directory into the go src directory
	symlink := exec.Command("ln", "-s", "-f", "/app/", "/go/src/"+gopackage)

	// redirect all output to the standard console
	symlink.Stdout = os.Stdout
	symlink.Stderr = os.Stderr

	err = symlink.Run()
	if err != nil {
		log.Fatalf("Command finished with error: %s\n", err)
	}

	// create a new file watcher utilizing inotify
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatalln("Error while watching for changes: %s\n", err)
	}

	// set the watcher path on the volume directly as symlinks are not followed
	// by inotify
	err = filepath.Walk("/app", watchIfDirectory(watcher))
	if err != nil {
		log.Fatalln(err)
	}

	// rebuild the package
	rebuildPackage(gopackage)

	for {
		select {
		case ev := <-watcher.Event:

			// only act on IN_MODIFY, IN_CREATE and IN_DELETE events
			if (ev.Mask&inotify.IN_MODIFY) == inotify.IN_MODIFY ||
				(ev.Mask&inotify.IN_CREATE) == inotify.IN_CREATE ||
				(ev.Mask&inotify.IN_DELETE) == inotify.IN_DELETE {

				// TODO: filter out all files that do not end in *.go or *.tmpl

				log.Printf("- %s\n", ev.Name)
				go rebuildPackage(gopackage)

			}

		case err := <-watcher.Error:
			log.Println("error:", err)
		}
	}

}

// rebuildPackage will rebuild the given package and restart the process
func rebuildPackage(gopackage string) {

	// build and install the package
	builder := exec.Command("go", "install", gopackage)

	// redirect all output to the standard console
	builder.Stdout = os.Stdout
	builder.Stderr = os.Stderr

	// run the build command and wait for it to exit
	err := builder.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not build the package", err)
		log.Fatalln("stopping process")
	}

	// TODO: get the package name to run the package

	// stop all previous running instances of the project
	kill := exec.Command("pkill", "-x", "test-project")

	// redirect all output to the standard console
	kill.Stdout = os.Stdout
	kill.Stderr = os.Stderr

	err = kill.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not stop running instances", err)
	}

	// TODO: make sure that only one process is running (i.e. stop the processes before)
	// explore use of pkill -f "auto-run"
	runner := exec.Command("test-project")

	// set the current directory to the gopackage
	os.Chdir("/go/src/" + gopackage)

	// redirect all output to the standard console
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr

	// run the program and wait for it to exit
	err = runner.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not run the package", err)
	}

}

// watchIfDirectory will return a function to watch all subdirectories of the given path
func watchIfDirectory(watcher *inotify.Watcher) filepath.WalkFunc {

	return func(path string, info os.FileInfo, err error) error {

		// skip everything that is not a directory
		if info.IsDir() == false {
			return nil
		}

		// skip some directories (i.e. vendor and versioning)
		excludeDirs := []string{"/vendor", "/node_modules", ".git", ".svn"}

		if containsAny(path, excludeDirs) == true {
			return filepath.SkipDir
		}

		// watch all other directories
		watcher.Watch(path)

		return err
	}

}
