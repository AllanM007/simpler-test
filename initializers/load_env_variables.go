package initializers

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {

	// Get the current directory of the project
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	err = godotenv.Load(filepath.Join(dir, ".env"))
	if err != nil {
		log.Println("error loading env file")
		log.Fatal(err)
	}
}
