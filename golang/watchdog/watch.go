package main

import (
	"golang.org/x/exp/inotify"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// watchForChanges will watch the /app volume for changes and rebuild the
// package as soon as changes occur
func watchForChanges(config Config) {

	packagePath := config.ProjectPath + config.Directory

	// rebuild and start the package
	restartPackage(packagePath, config.Arguments)

	// create a new file watcher utilizing inotify
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatalln("Error while watching for changes: %s\n", err)
	}

	// set the watcher path on the volume directly as symlinks are not followed
	// by inotify
	err = filepath.Walk("/app", initWatchlist(watcher, config.Ignore))
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case ev := <-watcher.Event:

			// rebuild the package if something was modified
			if (ev.Mask & inotify.IN_MODIFY) == inotify.IN_MODIFY {

				// rebuild and restart the package
				go restartPackage(packagePath, config.Arguments)

			} else if (ev.Mask & inotify.IN_CREATE) == inotify.IN_CREATE {

				// we need to add newly created directories as well
				addWatch(watcher, ev.Name)

				// rebuild and restart the package
				go restartPackage(packagePath, config.Arguments)

			} else if (ev.Mask & inotify.IN_DELETE) == inotify.IN_DELETE {

				// we need to add newly created directories as well
				removeWatch(watcher, ev.Name)

				// rebuild and restart the package
				go restartPackage(packagePath, config.Arguments)

			}

		case err := <-watcher.Error:
			log.Println("error:", err)
		}
	}

}

// initWatchlist will return a function to watch all subdirectories of the given path
func initWatchlist(watcher *inotify.Watcher, ignore []string) filepath.WalkFunc {

	// skip some directories by default (i.e. vendor and versioning)
	excludeDirs := []string{"/vendor", "/node_modules", ".git", ".svn"}

	// prefix all ignored directory with /app to create an absolute path
	ignorePaths := make([]string, len(ignore))
	for index, value := range ignore {
		ignorePaths[index] = "/app/" + value
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
		watcher.Watch(path)

		return err
	}

}

// addWatch will include newly created directories in the
// watchlist
func addWatch(watcher *inotify.Watcher, path string) {

	info, _ := os.Stat(path)
	if info.IsDir() {
		watcher.Watch(path)
	}

}

// removeWatch will remove the watcher for the given path
func removeWatch(watcher *inotify.Watcher, path string) {

	// note: this will return an error if the watch does not exist, but we
	// do not need to care about that
	watcher.RemoveWatch(path)
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
