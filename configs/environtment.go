package configs

import (
	"fmt"

	"gofit-api/models"

	"github.com/spf13/viper"
)

var (
	AppConfig models.Config
)

func LoadConfig() *models.Config {
	viper.SetConfigType("env")
	viper.SetConfigName("local")
	// viper.SetConfigName("dev")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return &AppConfig
}

func SetMidtransConfig()  {
	
}