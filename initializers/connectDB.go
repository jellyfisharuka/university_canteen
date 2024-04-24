package initializers

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func ConnectDb() {
	var err error
	db_connect := os.Getenv("db")
	DB, err = gorm.Open(postgres.Open(db_connect), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to DB")
	}
}
