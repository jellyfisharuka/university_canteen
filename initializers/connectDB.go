package initializers

import (
	"final_project/internal/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDb() {
	var err error
	db_connect := os.Getenv("db")
	DB, err = gorm.Open(postgres.Open(db_connect), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to DB")
	}
	DB.AutoMigrate(models.Users{})
}
