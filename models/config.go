package models

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDatabase string `mapstructure:"REDIS_DATABASE"`

	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPUsername string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`

	MidtransMerchantID string `mapstructure:"MIDTRANS_MERCHANT_ID"`
	MidtransClientKey  string `mapstructure:"MIDTRANS_CLIENT_KEY"`
	MidtransServerKey  string `mapstructure:"MIDTRANS_SERVER_KEY"`

	JWTKey  string `mapstructure:"JWT_KEY"`
	AppPort string `mapstructure:"APP_PORT"`
}
