package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// parseConfiguration will parse the necessary external information from the command line
// or the environment and return an error if the flag is not defined
func parseConfiguration() (Config, error) {

	// the package import path should be supplied via flag or environment variable
	config := Config{}

	// initialize a string for our ignore values
	var ignore string
	var arguments string
	var timeout string

	// parse additional information from the command line
	flag.StringVar(&config.Directory, "directory", defaultDirectory, "(optional) absolute path of the go module directory inside the docker container")
	flag.StringVar(&ignore, "ignore", "", "(optional) directories to ignore when watching for changes")
	flag.StringVar(&arguments, "args", "", "(optional) arguments to pass to the service on start")
	flag.StringVar(&config.Command, "cmd", "build", "(optional) use 'build' to auto restart the code, 'test' to automatically run 'go test', 'noop' to not run anything")
	flag.StringVar(&timeout, "timeout", "1500ms", "(optional) timeout to wait for further file changes until restart is triggered")

	flag.Parse()

	if config.Directory == "" {
		config.Directory = os.Getenv("DIRECTORY")
	}

	if ignore == "" {
		ignore = os.Getenv("IGNORE")
	}

	if arguments == "" {
		arguments = os.Getenv("ARGS")
	}

	if timeout == "500ms" {
		// allow overriding of the default timeout from environment
		envCommand := os.Getenv("TIMEOUT")
		if envCommand != "" {
			timeout = envCommand
		}
	}

	if config.Command == "build" {
		// allow overriding of the default command from environment
		envCommand := os.Getenv("CMD")
		if envCommand != "" {
			config.Command = envCommand
		}
	}

	switch config.Command {
	case "noop":
		fmt.Println("please log into the container to run any commands")

	case "build", "test":
		// nothing to do

	default:
		err := fmt.Errorf("the command '%s' is not defined. please use build or test", config.Command)
		return config, err
	}

	// ensure that the directory path starts with a slash
	if config.Directory != "" && strings.HasPrefix(config.Directory, "/") == false {
		config.Directory = "/" + config.Directory
	}

	if ignore != "" {
		config.Ignore = strings.Split(ignore, ",")

		for index, value := range config.Ignore {
			value = strings.TrimSpace(value)
			config.Ignore[index] = strings.TrimLeft(value, "/")
		}
	}

	if arguments != "" {
		config.Arguments = strings.Split(arguments, " ")
	}

	// parse the timeout duration
	var err error
	config.Timeout, err = time.ParseDuration(timeout)
	if err != nil {
		fmt.Printf("could not parse timeout duration: %s", timeout)
		config.Timeout = time.Millisecond * 500
	}

	return config, nil

}
