package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)
func LoadEnv() error {

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	file := ".env." + env

	if err := godotenv.Load(file); err != nil {
		return err
	}

	log.Println("Loaded environment:", env)
	return nil
}
