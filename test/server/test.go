package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/configor"
)

type Configuration struct {
	Message string
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func main() {

	config := Configuration{}
	err := configor.Load(&config, "config.yaml")
	if err != nil {
		fmt.Println("Configuration settings incomplete:", err.Error())
	}
	fmt.Printf("It works: %s\n", config.Message)
	fmt.Println("And now also with automatic updates..")

	taskid := uuid(10)

	for {
		fmt.Printf("%s: still running\n", taskid)
		time.Sleep(1 * time.Second)
	}

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
