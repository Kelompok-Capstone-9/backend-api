package main

import (
	"gofit-api/configs"
	assetsmanager "gofit-api/lib/assets_manager"
	"gofit-api/lib/scheduler"
	"gofit-api/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	err := assetsmanager.InitAssetsFile()
	if err != nil {
		panic(err)
	}

	configs.LoadConfig()
	err = configs.InitDB()
	if err != nil {
		panic(err)
	}

	err = configs.MigrateDB()
	if err != nil {
		panic(err)
	}

	// err = configs.SeedDB()
	// if err != nil {
	// 	panic(err)
	// }

	// err = configs.MigrateAndSeedDB()
	// if err != nil {
	// 	panic(err)
	// }

	e := echo.New()

	go scheduler.ScheduleMembershipActivityCheck()

	routes.InitRoute(e)

	e.Logger.Fatal(e.Start(":" + configs.AppConfig.AppPort))

}
