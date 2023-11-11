package main

import (
	"fmt"
	"log"
	// "github.com/Valgard/godotenv"
)

func main() {
	// dotenv := godotenv.New()
	// if err := dotenv.Load(".env"); err != nil {
	// 	panic(err)
	// }
	fmt.Println("Booting up...")
	store, err := NewPostgresStore()

	fmt.Println("Postgres Store Created...")
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Store initialized...")
	// port := os.Getenv("PORT")
	port := "3000"
	server := NewAPIServer(":"+port, store)
	fmt.Printf("Server running on port %s", port)
	server.Run()
}
