package main

import (
	"golang.org/x/exp/inotify"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// watchPackageDirectory will watch the /app volume for changes and rebuild the
// package as soon as changes occur
func watchPackageDirectory(gopackage string) {

	// create a new file watcher utilizing inotify
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatalln("Error while watching for changes: %s\n", err)
	}

	// set the watcher path on the volume directly as symlinks are not followed
	// by inotify
	err = filepath.Walk("/app", initWatchlist(watcher))
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case ev := <-watcher.Event:

			// rebuild the package if something was modified
			if (ev.Mask & inotify.IN_MODIFY) == inotify.IN_MODIFY {

				// TODO: filter out all files that do not end in *.go or *.tmpl

				// rebuild and restart the package
				go restartPackage(gopackage)

			} else if (ev.Mask & inotify.IN_CREATE) == inotify.IN_CREATE {

				// we need to add newly created directories as well
				addWatch(watcher, ev.Name)

				// rebuild and restart the package
				go restartPackage(gopackage)

			} else if (ev.Mask & inotify.IN_DELETE) == inotify.IN_DELETE {

				// we need to add newly created directories as well
				removeWatch(watcher, ev.Name)

				// rebuild and restart the package
				go restartPackage(gopackage)

			}

		case err := <-watcher.Error:
			log.Println("error:", err)
		}
	}

}

// initWatchlist will return a function to watch all subdirectories of the given path
func initWatchlist(watcher *inotify.Watcher) filepath.WalkFunc {

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
