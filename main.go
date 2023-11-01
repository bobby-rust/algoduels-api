package main

import (
	"fmt"
)

func main() {
	fmt.Println("Heoo")

	server := NewAPIServer(":3000")
	server.Run()
}	