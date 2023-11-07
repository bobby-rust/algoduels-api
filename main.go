package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Valgard/godotenv"
)

func main() {
	dotenv := godotenv.New()
	if err := dotenv.Load(".env"); err != nil {
		panic(err)
	}

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	server := NewAPIServer(":" + port, store)
	fmt.Printf("Server running on port %s", port)
	server.Run()
}
