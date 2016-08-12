package main

import (
	"flag"
	"fmt"
	"github.com/jinzhu/configor"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type Configuration struct {
	Message         string `required:"true"`
	PublicDirectory string `required:"true"`
	Port            string `required:"true"`
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func main() {

	// load some configuration to check vendoring support
	config := Configuration{}
	err := configor.Load(&config, "config.yaml")
	if err != nil {
		fmt.Println("Configuration settings incomplete:", err.Error())
	}

	fmt.Printf("Message: %s\n", config.Message)
	fmt.Printf("Host: %s\n", config.Port)
	fmt.Printf("PublicDirectory: %s\n", config.PublicDirectory)

	// try to parse arguments from the command line
	var testArgument string
	flag.StringVar(&testArgument, "test", "not defined", "please specify a test argument")
	flag.Parse()

	fmt.Printf("Flag test: %s\n", testArgument)

	// start a simple webserver serving the assets directory and providing a
	// simple api call
	http.HandleFunc("/api", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "UUID: "+uuid(10)+"\n")
	})

	// serve the index file on the main port
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, config.PublicDirectory+"/index.html")
	})

	fs := http.FileServer(http.Dir(config.PublicDirectory))
	http.Handle("/assets", fs)

	fmt.Println("Starting server on", config.Port)
	http.ListenAndServe(config.Port, nil)

}

func uuid(strlen int) string {

	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}
