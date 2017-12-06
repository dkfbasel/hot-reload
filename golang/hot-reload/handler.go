package main

import (
	"sync"
	"time"
)

// Handler is used to handle package restarts
type Handler struct {
	notifications    chan (bool)
	lastNotification time.Time
	rebuild          bool
	packagePath      string
	executable       string
	arguments        []string
	sync.RWMutex
}

func initHandler(packagePath string, executable string, arguments []string) *Handler {

	// initialize the handler struct, set rebuild to true to automatically build
	// the package on initialization
	handler := Handler{
		notifications: make(chan bool),
		rebuild:       true,
		packagePath:   packagePath,
		executable:    executable,
		arguments:     arguments,
	}

	go func() {
		for {
			select {

			// react to notifications from file watching
			case <-handler.notifications:
				handler.Lock()
				handler.rebuild = true
				handler.lastNotification = time.Now()
				handler.Unlock()

			// check every 100 milliseconds if at least 300 milliseconds have
			// elapsed since the last notification. if so, restart the build/test process
			case <-time.After(100 * time.Millisecond):
				handler.RLock()
				if handler.rebuild == true && time.Now().After(handler.lastNotification.Add(300*time.Millisecond)) {
					handler.rebuild = false
					go handler.restartPackage()
				}
				handler.RUnlock()
			}

		}
	}()

	return &handler

}
