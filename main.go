package main

import (
	"log"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	log.Println("Successfully connected to the database")

	if err := store.Init(); err != nil {
		log.Fatalf("Error initializing the database: %v", err)
	}
	log.Println("Successfully initialized the database")

	server := NewAPIServer(":8008", store)
	server.Run()
}
 