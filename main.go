package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	// db
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	err = store.Init()
	if err != nil {
		os.Exit(1)
	}
	
	server := NewAPIService(":3000", store)
	server.Run()
	fmt.Println("hello world!")
}
