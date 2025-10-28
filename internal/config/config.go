package config

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	SecretToken          string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf(".env не найден: %v", err)
	} else {
		log.Println(".env загружен")
	}

	runAddr := getEnv("RUN_ADDRESS", ":8080")
	dbURI := getEnv("DATABASE_URL", "")
	accrualAddr := getEnv("ACCRUAL_SYSTEM_ADDRESS", "http://localhost:8081")
	secretToken := getEnv("SECRET_TOKEN", "")

	flag.StringVar(&runAddr, "a", runAddr, "адрес сервера (RUN_ADDRESS)")
	flag.StringVar(&dbURI, "d", dbURI, "URI базы данных (DATABASE_URL)")
	flag.StringVar(&accrualAddr, "r", accrualAddr, "адрес системы начислений (ACCRUAL_SYSTEM_ADDRESS)")
	flag.Parse()

	return &Config{
		RunAddress:           runAddr,
		DatabaseURI:          dbURI,
		AccrualSystemAddress: accrualAddr,
		SecretToken:          secretToken,
	}
}

func getEnv(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	return value
}
