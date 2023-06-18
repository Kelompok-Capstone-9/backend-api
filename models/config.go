package models

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	RedisAddress string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDatabase string `mapstructure:"REDIS_DATABASE"`

	SMTPHost string `mapstructure:"SMTP_HOST"`
	SMTPUsername string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
}