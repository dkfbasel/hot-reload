package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/becheran/wildmatch-go"
	"github.com/fsnotify/fsnotify"
)

// watchForChanges will watch the /app volume for changes and rebuild the
// package as soon as changes occur
func watchForChanges(config Config, notify chan<- bool) {

	// create a new file watcher utilizing inotify
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error while watching for changes: %s\n", err)
	}
	defer watcher.Close() // nolint: errcheck

	// prefix all ignored directory with /app to create an absolute path
	ignoreList := make([]string, len(config.Ignore))
	for index, value := range config.Ignore {
		ignoreList[index] = path.Join(config.Directory, value)
	}

	// initialize a map of all ignored files, since we need to cancel the
	// notification on this files, if they are inside of a watched folder
	ignoredFiles := make(map[string]struct{})

	// set the watcher path on the volume directly as symlinks are not followed
	// by inotify
	fmt.Printf("%s SETUP\n----------------------------\n", time.Now().Format("2006/01/02 15:04:05"))
	err = filepath.WalkDir(config.Directory, initWatchlist(watcher, config.Directory,
		ignoreList, ignoredFiles))
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case event := <-watcher.Events:

			// rebuild the package if something was modified
			if event.Op&fsnotify.Write == fsnotify.Write {
				// ev.Mask & inotify.IN_MODIFY) == inotify.IN_MODIFY

				// ignore all items that are in our ignoredFiles list
				if _, ok := ignoredFiles[event.Name]; ok {
					continue
				}

				// rebuild and restart the package
				notify <- true

			} else if event.Op&fsnotify.Create == fsnotify.Create {
				// if (ev.Mask & inotify.IN_CREATE) == inotify.IN_CREATE

				// check if the new item should be added to our ignore list
				if matchesAny(event.Name, ignoreList) {
					// ignore the path if there should be any notification, this will
					// also include new directories, but does not really concern us
					ignoredFiles[event.Name] = struct{}{}
					continue
				}

				// we need to add newly created directories as well
				addWatch(watcher, event.Name)

				// rebuild and restart the package
				notify <- true

			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				// (ev.Mask & inotify.IN_DELETE) == inotify.IN_DELETE

				// we need to add newly created directories as well
				removeWatch(watcher, event.Name)

				// rebuild and restart the package
				notify <- true

			}

		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}

}

// initWatchlist will return a function to watch all subdirectories of the given path
func initWatchlist(watcher *fsnotify.Watcher, directory string,
	ignoreList []string, ignoredFiles map[string]struct{}) fs.WalkDirFunc {

	// skip some directories by default (i.e. vendor and versioning)
	excludeDirs := []string{"/vendor", "/node_modules", ".git", ".svn"}

	// go through all directories
	return func(path string, info fs.DirEntry, err error) error {

		if containsAny(path, excludeDirs) {
			fmt.Printf("watch: ignore path %s\n", path)
			watcher.Remove(path) // nolint: errcheck
			return filepath.SkipDir
		}

		// ignore all directories that have been specified to be skipped
		if matchesAny(path, ignoreList) {
			fmt.Printf("watch: ignore path %s\n", path)
			watcher.Remove(path) // nolint: errcheck
			if info.IsDir() {
				return filepath.SkipDir
			}
			// add the file to the list of ignored files, to cancel watcher notifications
			// effectively later on
			ignoredFiles[path] = struct{}{}
			return nil
		}

		// do not add watchers on specific files
		if !info.IsDir() {
			return nil
		}

		// watch all other directories
		// fmt.Printf("watch: add path %s\n", path)
		watcher.Add(path) // nolint: errcheck

		return err
	}

}

// addWatch will include newly created directories in the watchlist
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
		if strings.Contains(source, element) {
			return true
		}
	}
	return false
}

// matchesAny will check whether the source matches any of the given patterns
func matchesAny(path string, patterns []string) bool {

	for _, pattern := range patterns {
		// use wildmatch-library for simple pattern matching
		matcher := wildmatch.NewWildMatch(pattern)
		if matcher.IsMatch(path) {
			return true
		}
	}

	return false

}
