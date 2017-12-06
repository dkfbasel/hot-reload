package main

import (
	"fmt"
	"os"
	"strings"
)

// default directory to watch is the /app directory
var directoryToWatch = "/app"

func main() {

	// parse the gopackage from the command line or an environment variable.
	// note: a tmp directory will be used if no project path is supplied, this
	// might break some import statements that point to the project directory
	config := parseConfiguration()

	// print the information that was parsed from the flags
	fmt.Printf("PACKAGE:    %s\n", config.ProjectPath)

	if config.Directory != "" {
		fmt.Printf("DIRECTORY:  %s\n", config.Directory)
	}

	if len(config.Ignore) > 0 {
		fmt.Printf("IGNORE:     %s\n", strings.Join(config.Ignore, ", "))
	}

	if len(config.Arguments) > 0 {
		fmt.Printf("ARGUMENTS:  %s\n", strings.Join(config.Arguments, " "))
	}

	if config.ProjectPath == tmpProjectPath {
		fmt.Println("please note that import paths in the project directory will probably not work as intended")
	}

	// check if package directory exists. if not, create a symlink from the /app
	// directory. note: this will not work with goconvey (goconvey cannot follow symlinks)
	if _, err := os.Stat("/go/src/" + config.ProjectPath); os.IsNotExist(err) {
		createSymlinkForPackage(config)
	} else {
		// watch the package directory directly
		directoryToWatch = "/go/src/" + config.ProjectPath
	}

	switch config.Command {
	case "build", "test":
		// watch the supplied directory for changes and rebuild and rerun the package
		watchForChanges(config.Command, directoryToWatch, config)

	case "goconvey":
		// inform user that ignore directories must currently be specified differently
		if len(config.Ignore) > 0 {
			fmt.Println("please use the argument -excludeDirs from goconvey to exclude directories")
		}

		// start goconvey (will watch directories automatically)
		runGoconvey(config.ProjectPath, config.Arguments)

	case "noop":
		fmt.Println("please log into the container to run any commands")

	default:
		fmt.Printf("the command '%s' is not defined. please use build, goconvey or test", config.Command)
	}

}
