package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const tmpProjectPath = "development.com/tmp"

// Config will contain all flags than can
type Config struct {
	ProjectPath string   // the base path to the project
	Directory   string   // the subdirectory containing the current package
	Ignore      []string // directories to ignore when watching for changes
	Arguments   []string // arguments to pass to the service
}

func main() {

	// parse the gopackage from the command line or an environment variable.
	// note: a tmp directory will be used if no project path is supplied, this
	// might break some import statements that point to the project directory
	config := parseConfiguration()

	// print the information that was parsed from the flags
	fmt.Printf("PACKAGE:    %s\n", config.ProjectPath)
	fmt.Printf("DIRECTORY:  %s\n", config.Directory)
	fmt.Printf("IGNORE:     %s\n", strings.Join(config.Ignore, ", "))
	fmt.Printf("ARGUMENTS:  %s\n", strings.Join(config.Arguments, " "))

	if config.ProjectPath == tmpProjectPath {
		fmt.Printf("please note that import paths in the project directory will probably not work as intended")
	}

	// create a symlink for the package to allow for compilation
	createSymlinkForPackage(config.ProjectPath)

	// watch the supplied directory for changes and rebuild and rerun the package
	watchForChanges(config)

}

// parseConfiguration will parse the necessary external information from the command line
// or the environment and return an error if the flag is not defined
func parseConfiguration() Config {

	// the package import path should be supplied via flag or environment variable
	config := Config{}

	// initialize a string for our ignore values
	var ignore string
	var arguments string

	// parse additional information from the command line
	flag.StringVar(&config.ProjectPath, "project", "", "the path of the project relative to the go source directory")
	flag.StringVar(&config.Directory, "directory", "", "(optional) relative path to the directory of the current go package to watch")
	flag.StringVar(&ignore, "ignore", "", "(optional) directories to ignore when watching for changes")
	flag.StringVar(&arguments, "args", "", "(optional) arguments to pass to the service on start")
	flag.Parse()

	// try to parse the data from the environment if not supplied via flag
	if config.ProjectPath == "" {
		config.ProjectPath = os.Getenv("PROJECT")
	}

	if config.Directory == "" {
		config.Directory = os.Getenv("DIRECTORY")
	}

	if ignore == "" {
		ignore = os.Getenv("IGNORE")
	}

	if arguments == "" {
		arguments = os.Getenv("ARGUMENTS")
	}

	// set the project path to tmp if not specified -> this will break package
	// imports that are contained in the same project directory
	if config.ProjectPath == "" {
		config.ProjectPath = tmpProjectPath
	}

	// ensure that the subdirectory starts with a slash
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

	return config

}
