package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/AllanM007/simpler-test/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() *gorm.DB {
	var err error
	DB, err = ConnectDB()
	if err != nil {
		log.Fatalf("Connection to database failed")
	}
	fmt.Printf("Database connection successfully established")
	return DB
}

func ConnectDB() (*gorm.DB, error) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password='%s' dbname=%s port=%s TimeZone=Africa/Nairobi", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to database : error=%w", err)
	}

	return db, nil
}

func MigrateDB() error {
	if DB == nil {
		return fmt.Errorf("no database connection available")
	}
	return DB.AutoMigrate(&models.Product{})
}
