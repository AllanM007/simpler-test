package initializers

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// Test function for loading .env file
func TestLoadEnvVariables(t *testing.T) {
	// Load the .env file
	// LoadEnvVariables()
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Test that the environment variables are loaded correctly
	expectedVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}

	for _, envVar := range expectedVars {
		value := os.Getenv(envVar)
		if value == "" {
			t.Errorf("Expected environment variable %s to be set, but it is empty", envVar)
		}
	}
}
