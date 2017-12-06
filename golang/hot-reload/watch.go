package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// getExecutableName will get the last part of the package path that is used
// to name the go executable
func getExecutableName(packagePath string) string {
	parts := strings.Split(packagePath, "/")
	return parts[len(parts)-1]
}

// watchForChanges will watch the /app volume for changes and rebuild the
// package as soon as changes occur
func watchForChanges(command string, directory string, config Config) {

	packagePath := config.ProjectPath + config.Directory

	// get the package name to run the package
	executable := getExecutableName(packagePath)

	// use go test for testing
	if command == "test" {
		executable = "go test"
	}

	// initialize a change handler to react to changes
	changeHandler := initHandler(packagePath, executable, config.Arguments)

	// create a new file watcher utilizing inotify
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln("Error while watching for changes: %s\n", err)
	}
	defer watcher.Close() // nolint: errcheck

	// set the watcher path on the volume directly as symlinks are not followed
	// by inotify
	err = filepath.Walk(directory, initWatchlist(watcher, directory, config.Ignore))
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case event := <-watcher.Events:

			// rebuild the package if something was modified
			if event.Op&fsnotify.Write == fsnotify.Write {
				// ev.Mask & inotify.IN_MODIFY) == inotify.IN_MODIFY

				// rebuild and restart the package
				changeHandler.notifications <- true

			} else if event.Op&fsnotify.Create == fsnotify.Create {
				// if (ev.Mask & inotify.IN_CREATE) == inotify.IN_CREATE

				// we need to add newly created directories as well
				addWatch(watcher, event.Name)

				// rebuild and restart the package
				changeHandler.notifications <- true

			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				// (ev.Mask & inotify.IN_DELETE) == inotify.IN_DELETE

				// we need to add newly created directories as well
				removeWatch(watcher, event.Name)

				// rebuild and restart the package
				changeHandler.notifications <- true

			}

		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}

}

// initWatchlist will return a function to watch all subdirectories of the given path
func initWatchlist(watcher *fsnotify.Watcher, directory string, ignore []string) filepath.WalkFunc {

	// skip some directories by default (i.e. vendor and versioning)
	excludeDirs := []string{"/vendor", "/node_modules", ".git", ".svn"}

	// prefix all ignored directory with /app to create an absolute path
	ignorePaths := make([]string, len(ignore))
	for index, value := range ignore {
		ignorePaths[index] = path.Join(directory, value)
	}

	// go through all directories
	return func(path string, info os.FileInfo, err error) error {

		// skip everything that is not a directory
		if info.IsDir() == false {
			return nil
		}

		if containsAny(path, excludeDirs) == true {
			return filepath.SkipDir
		}

		// ignore all directories that have been specified to be skipped
		if equalsAny(path, ignorePaths) {
			return filepath.SkipDir
		}

		// watch all other directories
		watcher.Add(path) // nolint: errcheck

		return err
	}

}

// addWatch will include newly created directories in the
// watchlist
func addWatch(watcher *fsnotify.Watcher, path string) {

	info, err := os.Stat(path)
	if err != nil {
		return
	}

	if info.IsDir() {
		watcher.Add(path) // nolint: errcheck
	}

}

// removeWatch will remove the watcher for the given path
func removeWatch(watcher *fsnotify.Watcher, path string) {

	// note: this will return an error if the watch does not exist, but we
	// do not need to care about that
	watcher.Remove(path) // nolint: errcheck
}

// containsAny will check whether any of the matches is part of the given string
func containsAny(source string, matches []string) bool {

	for _, element := range matches {
		if strings.Contains(source, element) == true {
			return true
		}
	}

	return false
}

// equalsAny will check whether the source equals any of the given strings
func equalsAny(source string, matches []string) bool {

	for _, element := range matches {
		if source == element {
			return true
		}
	}

	return false

}
