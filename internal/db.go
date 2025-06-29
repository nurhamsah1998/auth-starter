package internal

import (
	"fmt"
	"log"

	"github.com/nurhamsah1998/auth-starter/config"
	"github.com/nurhamsah1998/auth-starter/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DbGormInit() {
	cfg := config.LoadConfig()
	dbENV := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.DbName, cfg.Port, cfg.SSLMode)
	db, err := gorm.Open(postgres.Open(dbENV), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	DB = db

	errMigrate := db.AutoMigrate(&model.User{}, &model.Profile{})
	if errMigrate != nil {
		log.Fatal("Migrate Failed", err)
	}
	println("Successfully connect to DB")
}
