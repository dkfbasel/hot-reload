package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	// parse the gopackage from the command line or an environment variable
	gopackage, err := parseFlag("gopackage", "import path of the go package to watch")

	if err != nil {
		fmt.Println("ERROR: A go package import path is required.\n-e \"GOPACKAGE=bitbucket.com/dkfbasel/test\"\n\n")
		return
	}

	// print the package path for the user
	fmt.Printf("PACKAGE:\n%s\n", gopackage)

	// create a symlink for the package to allow for compilation
	createSymlinkForPackage(gopackage)

	// rebuild and start the package
	restartPackage(gopackage)

	// watch the package directory for changes
	watchPackageDirectory(gopackage)
}

// parseFlag will parse the given flag from the command line or the environment
// and return an error if the flag is not defined
func parseFlag(flagName, info string) (string, error) {

	// the package import path should be supplied via flag or environment variable
	var content string
	flag.StringVar(&content, flagName, "", info)
	flag.Parse()

	// try to parse the data from the environment if not supplied via flag
	if content == "" || content == "undefined" {
		content = os.Getenv(strings.ToUpper(flagName))
	}

	if content == "" || content == "undefined" {
		return "", errors.New("flag not supplied: " + flagName)
	}

	return content, nil

}
