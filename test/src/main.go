package main

import (
	"fmt"
	"net/http"

	"github.com/gookit/color"
	"github.com/kelseyhightower/envconfig"
)

type Specification struct {
	Debug bool
	Port  int `default:"80"`
}

func main() {

	var config Specification
	err := envconfig.Process("myapp", &config)
	if err != nil {
		fmt.Printf("[ERROR] could not parse config: %+v\n", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello whats up")
	})

	color.Cyan.Println("We are running inside the container to make this work")

	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.Port), nil)
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
	}
	fmt.Println("something")
}
