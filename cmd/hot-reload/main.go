package main

import (
	"log"
	"time"
)

// default directory to watch is the /app directory
var defaultDirectory = "/app"

// Config contains all flags than can be passed to the utility
type Config struct {
	Directory string        // the main directory of the project
	Command   string        // the command to use for watching
	Ignore    []string      // directories to ignore when watching for changes
	Arguments []string      // arguments to pass to the service
	Timeout   time.Duration // timeout as time string (i.e 300ms)
	Proxy     string        // address of the app which should be proxied
}

func main() {

	// parse the configration from the command line or environment variables
	config, err := parseConfiguration()
	if err != nil {
		log.Printf("[ERROR] CONFIGURATION ISSUE:\n%+v", err)
		return
	}

	// initialize a notify channel to handle any file system changes
	notifyChan := make(chan bool)
	abortNotify := make(chan bool)

	// wait for input on the notify channel
	go func() {
		// handle all notifications
		for range notifyChan {

			// abort any waiting routine using a non blocking send operation
			// which will only trigger if there is currently an open receiver
			select {
			case abortNotify <- true:
			default:
			}

			go func() {
				// wait some time before the run command is started, unless it
				// is aborted beforehand through a new file change action
				select {
				case <-time.After(config.Timeout):
					switch config.Command {
					case "build":
						runBuild(config)
						broadcast("reload")

					case "test":
						runTest(config)
					}

				case <-abortNotify:
					// abort running the command
				}
			}()

		}
	}()

	// initialize the first build
	notifyChan <- true

	// run a proxy web server to handle hot reload requests
	go runHttpServer(config)

	// watch the supplied directory for changes
	watchForChanges(config, notifyChan)

}
