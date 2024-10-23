package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/AllanM007/simpler-test/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password='%s' dbname=%s port=%s TimeZone=Africa/Nairobi", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error connecting to database : error=%v", err)
	}

	fmt.Printf("Database connection successfully established")
	return db, nil
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Product{},
	)
}
