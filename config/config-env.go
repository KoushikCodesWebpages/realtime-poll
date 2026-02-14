package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func LoadEnv() error {

	// // 1️⃣ If APP_ENV is provided externally → respect it
	// env := os.Getenv("APP_ENV")

	// if env != "" {
	// 	file := "."+ env + ".env." 
	// 	if err := godotenv.Load(file); err != nil {
	// 		return err
	// 	}
	// 	log.Println("Loaded environment from:", file)
	// 	return nil
	// }

	// 2️⃣ Otherwise auto-pick file
	if exists(".dev.env") {
		if err := godotenv.Load(".dev.env"); err != nil {
			return err
		}
		log.Println("Loaded environment: dev (auto)")
		return nil
	}

	if exists(".prod.env") {
		if err := godotenv.Load(".prod.env"); err != nil {
			return err
		}
		log.Println("Loaded environment: prod (auto)")
		return nil
	}

	log.Fatal("No environment file found (.env.dev or .env.prod)")
	return nil
}
