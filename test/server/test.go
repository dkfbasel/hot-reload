package main

import (
	"fmt"
	"github.com/jinzhu/configor"
)

type Configuration struct {
	Message string
}

func main() {

	config := Configuration{}
	err := configor.Load(&config, "config.yaml")
	if err != nil {
		fmt.Println("Configuration settings incomplete:", err.Error())
	}
	fmt.Printf("It works: %s\n", config.Message)
	fmt.Println("And now also with automatic updates..")
}
