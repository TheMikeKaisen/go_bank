package main

import "fmt"

func main() {
	server := NewAPIService(":8000")
	server.Run()
	fmt.Println("hello world!")
}
