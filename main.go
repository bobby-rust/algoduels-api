package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	fmt.Println("Booting up...")
	store, err := NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Postgres Store Created...")

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	test(store)

	fmt.Println("Store initialized...")
	port := os.Getenv("PORT")
	server := NewAPIServer(":"+port, store)
	server.Run()

}
