package configs

import (
	"fmt"
	"gofit-api/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		AppConfig.DBUsername,
		AppConfig.DBPassword,
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBName,
	)
	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	return err
}

func MigrateDB() error {
	return DB.AutoMigrate(
		models.User{},
		models.Class{},
		models.Location{},
	)
}

func SeedDB() error {
	var (
		admin = models.User{
			Name:     "GoFit Administrator",
			Email:    "gofit@gofit.com",
			Password: "gofitadmin123",
			Gender:   "pria",
			Height:   158,
			Weight:   60,
			IsAdmin:  true,
		}

		offlineJakarta = models.Location{
			Name:      "Offline",
			City:      "Depok",
			Latitude:  "1572619562112",
			Longitude: "1527129572712",
			ClassID:   1,
		}

		offlineClass = models.Class{
			Name:        "Cardio Class",
			Description: "Kelas kebugaran untuk mengurangi lemak dalam tubuh",
			ClassType:   models.Offline,
			StartedAt:   time.Now(),
			Location:    models.Location{ID: 1},
		}
	)
	admin.HashingPassword(&models.CustomError{})
	err := DB.FirstOrCreate(&admin).Error
	if err != nil {
		return err
	}

	err = DB.FirstOrCreate(&offlineClass).Error
	if err != nil {
		return err
	}

	err = DB.FirstOrCreate(&offlineJakarta).Error
	return err
}
