package initializers

import (
	"fmt"
	"os"

	"github.com/AllanM007/simpler-test/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() *gorm.DB {
	Db = connectDB()
	return Db
}

func connectDB() *gorm.DB {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password='%s' dbname=%s port=%s TimeZone=Africa/Nairobi", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error connecting to database : error=%v", err)
		return nil
	}

	db.AutoMigrate(&models.Product{})

	return db
}
