package infrastructure

import (
	"os"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	WalletURL   string
	WalletToken string
	JWTSecret   string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		WalletURL:   os.Getenv("WALLET_URL"),
		WalletToken: os.Getenv("WALLET_TOKEN"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}
}
