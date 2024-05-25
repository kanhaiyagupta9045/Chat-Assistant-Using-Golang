package databases

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	_ "test/config"
	"test/models"
)

var db *gorm.DB

func InitDatabase() error {

	dsn := os.Getenv("dsn")
	if dsn == "" {
		log.Fatalf("DATABASE_DSN not set in environment")
		return fmt.Errorf("DATABASE_DSN not set in environment")
	}
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.File{}); err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err = sqlDB.Ping(); err != nil {
		return err
	}

	fmt.Println("Connected to database successfully")
	return nil
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatalf("Database not initialized. Please call InitDatabase first.")
	}
	return db
}
