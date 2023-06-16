package main

import (
	"gofit-api/configs"
	"gofit-api/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	configs.LoadConfig()
	err := configs.InitDB()
	if err != nil {
		panic(err)
	}

	err = configs.MigrateDB()
	if err != nil {
		panic(err)
	}

	// err = configs.MigrateAndSeedDB()
	// if err != nil {
	// 	panic(err)
	// }

	e := echo.New()
	routes.InitRoute(e)

	e.Logger.Fatal(e.Start(":" + configs.AppConfig.AppPort))
}
