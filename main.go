package main

import (
	"gofit-api/configs"
	"gofit-api/lib/scheduler"
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

	err = configs.SeedDB()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	go scheduler.ScheduleMembershipActivityCheck()

	routes.InitRoute(e)

	e.Logger.Fatal(e.Start(":" + configs.AppConfig.AppPort))

}
