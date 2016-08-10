package main

import (
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

	fmt.Printf("Loaded data from configuration: %s\n", config.Message)

	// taskid := uuid(10)
	//
	// for {
	// 	fmt.Printf("%s: still running\n", taskid)
	// 	time.Sleep(1 * time.Second)
	// }

	// start a simple webserver serving the assets directory and providing a
	// simple api call
	http.HandleFunc("/api", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello there: the uuid is: "+uuid(10)+"\n")
	})

	fs := http.FileServer(http.Dir(config.PublicDirectory))
	http.Handle("/", fs)

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
