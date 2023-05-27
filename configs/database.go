package configs

import (
	"fmt"
	"gofit-api/models"

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
	)
}

func SeedDB() error {
	var (
		admin = models.User{
			Name:     "M Fikri Ramadhan",
			Email:    "fikri@gmail.com",
			Password: "123",
			Gender:   "pria",
			Height:   158,
			Weight:   60,
			IsAdmin:  true,
		}
	)
	admin.HashingPassword(&models.CustomError{})
	err := DB.FirstOrCreate(&admin).Error

	return err
}
