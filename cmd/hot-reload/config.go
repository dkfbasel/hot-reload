package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const defaultTimeout = "800ms"
const defaultProxyPort = "3333"

// parseConfiguration will parse the necessary external information from the command line
// or the environment and return an error if the flag is not defined
func parseConfiguration() (Config, error) {

	// the package import path should be supplied via flag or environment variable
	config := Config{}

	// initialize a string for our ignore values
	var ignore string
	var arguments string
	var timeout string
	var proxytarget string
	var proxyport string

	// parse additional information from the command line
	flag.StringVar(&config.Directory, "directory", defaultDirectory, "(optional) absolute path of the go module directory inside the docker container")
	flag.StringVar(&ignore, "ignore", "", "(optional) directories and files to ignore when watching for changes")
	flag.StringVar(&arguments, "args", "", "(optional) arguments to pass to the service on start")
	flag.StringVar(&config.Command, "cmd", "build", "(optional) use 'build' to auto restart the code, 'test' to automatically run 'go test', 'noop' to not run anything")
	flag.StringVar(&timeout, "timeout", defaultTimeout, "(optional) timeout to wait for further file changes until restart is triggered")
	flag.StringVar(&proxytarget, "proxytarget", "", "(optional) address of the app which should be proxied. no proxy is used if left empty")
	flag.StringVar(&proxyport, "proxyport", defaultProxyPort, "(optional) port to run the proxy server on")

	flag.Parse()

	// check the environment for the directory flag
	if config.Directory == "" || config.Directory == defaultDirectory {
		envDir := os.Getenv("DIRECTORY")
		if envDir != "" {
			config.Directory = envDir
		}
	}

	if ignore == "" {
		ignore = os.Getenv("IGNORE")
	}

	if arguments == "" {
		arguments = os.Getenv("ARGS")
	}

	if timeout == defaultTimeout {
		// allow overriding of the default timeout from environment
		envCommand := os.Getenv("TIMEOUT")
		if envCommand != "" {
			timeout = envCommand
		}
	}

	if proxytarget == "" {
		// allow overriding of the default proxy from environment
		envProxy := os.Getenv("PROXYTARGET")
		if envProxy != "" {
			proxytarget = envProxy
		}
	}

	if proxyport == defaultProxyPort {
		// allow overriding of the default proxy port from environment
		envProxyPort := os.Getenv("PROXYPORT")
		if envProxyPort != "" {
			proxyport = envProxyPort
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
	if config.Directory != "" && !strings.HasPrefix(config.Directory, "/") {
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

	if proxytarget != "" {
		if !strings.HasPrefix(proxytarget, "http://") && !strings.HasPrefix(proxytarget, "https://") {
			proxytarget = "http://" + proxytarget
		}
	}
	config.ProxyTarget = proxytarget

	if proxyport != "" {
		config.ProxyPort = proxyport
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
