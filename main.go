package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/innv8/testcontainers/logic"
	"github.com/joho/godotenv"
)

func main() {
	var db *sql.DB
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env because %v", err)
	}

	db, err = logic.DBConnect(os.Getenv("DB_URI"))
	if err != nil {
		log.Fatalf("failed to connect to db because %v", err)
	}

	log.Printf("1. Create a new user")
	id, err := logic.CreateUser(db, "Sam", 50)
	if err != nil {
		log.Fatalf("Failed to create user because %v", err)
	}
	log.Printf("Created user : %d", id)

	log.Printf("2. Select User")
	user, err := logic.GetUser(db, id)
	if err != nil {
		log.Fatalf("Failed to get user %d because %v", id, err)
	}
	log.Printf("Selected user %d as %+v", id, user)

	log.Printf("3. Update User")
	err = logic.UpdateUser(db, id, "John", 43)
	if err != nil {
		log.Fatalf("Failed to update user %d because %v", id, err)
	}
	log.Printf("Updated user successfully")

	log.Printf("4. Delete User")
	err = logic.DeleteUser(db, id)
	if err != nil {
		log.Fatalf("Failed to delete user %d because %v", id, err)
	}
	log.Printf("Deleted user")
}
