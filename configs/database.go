package configs

import (
	"errors"
	"fmt"
	"gofit-api/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

var PaymentMethods = []models.PaymentMethod{
	{
		Name: "gopay",
	},
	{
		Name: "credit card",
	},
	{
		Name: "shoope pay",
	},
}

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
	err := DB.AutoMigrate(
		models.User{},
		models.Membership{},
		models.Plan{},
		models.Location{},
		models.Class{},
		models.ClassPackage{},
		models.ClassTicket{},
		models.PaymentMethod{},
		models.Transaction{},
    models.Healthtip{},
	)
	if err != nil {
		return err
	}

	for _, payment := range PaymentMethods {
		err = DB.FirstOrCreate(&payment).Error
	}

	return err
}

func MigrateAndSeedDB() error {
	// drop table if exists
	err := DB.Migrator().DropTable(models.User{}, models.Location{}, models.Class{}, models.ClassPackage{}, models.ClassTicket{})
	if err != nil {
		return errors.New("fail to drop table")
	}

	// migrate table
	err = DB.AutoMigrate(models.User{}, models.Location{}, models.Class{}, models.ClassPackage{}, models.ClassTicket{})
	if err != nil {
		return errors.New("fail to migrate")
	}

	var (
		users = []models.User{
			{
				Name:     "GoFit Administrator",
				Email:    "admin@gofit.com",
				Password: "gofitadmin123",
				Gender:   models.Pria,
				Height:   158,
				Weight:   60,
				IsAdmin:  true,
			},
			{
				Name:     "Edward Halley",
				Email:    "halley@gmail.com",
				Password: "halley123",
				Gender:   models.Pria,
				Height:   158,
				Weight:   60,
				IsAdmin:  false,
			},
			{
				Name:     "Katarina Snow",
				Email:    "katarina@gmail.com",
				Password: "katarina123",
				Gender:   models.Wanita,
				Height:   158,
				Weight:   60,
				IsAdmin:  false,
			},
		}

		location = models.Location{
			Name:      "GoFit Gym Depok",
			Address:   "Jl. Margonda Raya No.151",
			City:      "Depok",
			Latitude:  "1572619562112",
			Longitude: "1527129572712",
		}

		locationID uint = 1
		classes         = []models.Class{
			{
				Name:        "Cardio Class",
				Description: "Kelas kebugaran untuk mengurangi lemak dalam tubuh",
				ClassType:   models.Offline,
				StartedAt:   time.Now(),
				LocationID:  &locationID,
			},
			{
				Name:        "Yoga Class",
				Description: "Kelas Online kebugaran untuk mengurangi lemak dalam tubuh",
				ClassType:   models.Online,
				Link:        "https://zoom.us/yoga-room",
				StartedAt:   time.Now(),
			},
		}

		classPackages = []models.ClassPackage{
			{
				Period:  models.Daily,
				Price:   50000,
				ClassID: 1,
			},
			{
				Period:  models.Weekly,
				Price:   200000,
				ClassID: 1,
			},
			{
				Period:  models.Monthly,
				Price:   350000,
				ClassID: 1,
			},
			{
				Period:  models.Daily,
				Price:   30000,
				ClassID: 2,
			},
			{
				Period:  models.Weekly,
				Price:   150000,
				ClassID: 2,
			},
			{
				Period:  models.Monthly,
				Price:   250000,
				ClassID: 2,
			},
		}

		classTickets = []models.ClassTicket{
			{
				UserID:         2,
				ClassPackageID: 2,
				Status:         models.Booked,
			},
			{
				UserID:         3,
				ClassPackageID: 5,
				Status:         models.Pending,
			},
			{
				UserID:         2,
				ClassPackageID: 2,
				Status:         models.Cancelled,
			},
		}
	)

	for key := range users {
		users[key].HashingPassword(&models.CustomError{})
	}

	err = DB.Create(&users).Error
	if err != nil {
		return err
	}

	err = DB.Create(&location).Error
	if err != nil {
		return err
	}

	err = DB.Create(&classes).Error
	if err != nil {
		return err
	}

	err = DB.Create(&classPackages).Error
	if err != nil {
		return err
	}

	err = DB.Create(&classTickets).Error
	if err != nil {
		return err
	}

	return nil
}
