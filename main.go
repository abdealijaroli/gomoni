package main

import (
	"log"
)

// func seedAccount(store Storage) error {
// 	acc, err := GenerateNewAccount("John", "Doe", "john@doe.com", "verystrongpasswordhehe")
// 	if err != nil {
// 		return err
// 	}
// 	return store.CreateAccount(acc)
// }

// func dropDatabase(store Storage) error {
// 	return store.DropTable()
// }

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

	// if err := dropDatabase(store); err != nil {
	// 	log.Fatalf("Error dropping the database: %v", err)
	// }
	// log.Println("Successfully dropped the database")

	// if err := seedAccount(store); err != nil {
	// 	log.Fatalf("Error seeding the database: %v", err)
	// }
	// log.Println("Successfully seeded the database")

	server := NewAPIServer(":8008", store)
	server.Run()
}
