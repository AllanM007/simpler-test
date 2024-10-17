package tests

import (
	"log"
	"testing"

	"github.com/AllanM007/simpler-test/initializers"
	"github.com/joho/godotenv"
)

func TestDBConnection(t *testing.T) {
	// Load .env file before running the db connection test
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	db, err := initializers.ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	var version string
	err = db.Raw("SELECT version();").Scan(&version).Error
	if err != nil {
		t.Fatalf("Failed to execute version statement query: %v", err)
	}

	log.Printf("Succesfully connected to database. Version: %s", version)
}
