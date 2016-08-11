package main

import (
	"fmt"
	"github.com/jinzhu/configor"
	r "gopkg.in/dancannon/gorethink.v2"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type Configuration struct {
	Message         string `required:"true"`
	PublicDirectory string `required:"true"`
	Port            string `required:"true"`
	RethinkDB       struct {
		Address  string `required:"true"`
		Database string `required:"true"`
	}
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

	// start a simple webserver serving the assets directory and providing a
	// simple api call
	http.HandleFunc("/api", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "UUID: "+uuid(10)+"\n")
	})

	// connect to the rethinkdb database
	session, err := r.Connect(r.ConnectOpts{
		Address:  config.RethinkDB.Address,
		Database: config.RethinkDB.Database,
	})

	if err != nil {
		fmt.Println("could not connect to the rethinkdb service")
	}

	// test if access to a different service will work as well. i.e. the rethinkdb
	// database service
	http.HandleFunc("/api/db", func(w http.ResponseWriter, req *http.Request) {
		request, err := r.Expr("Hello World").Run(session)
		if err != nil {
			fmt.Println("could not run rethinkdb command\n", err)
		}

		var response string
		err = request.One(&response)

		if err != nil {
			fmt.Println("could not fetch result from rethinkdb\n", err)
		}

		io.WriteString(w, "RethinkDB: "+response)

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
