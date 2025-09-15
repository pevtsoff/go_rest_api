package main

import (
	"fmt"
	"log"
	"os"
	"rest_api/config"
	"rest_api/models"
)

func printEnvVars() {
	envVars := os.Environ()

	for _, envVar := range envVars {
		fmt.Println(envVar)
	}
}

func init() {
	config.LoadEnvVars()
	printEnvVars()
	config.ConnectToDB()
}

func main() {
	if config.DB != nil {
		log.Println("Starting database migration...")

		// Auto-migrate all your models
		err := config.DB.AutoMigrate(&models.Post{}, &models.User{})
		if err != nil {
			log.Fatal("Migration failed:", err)
		}

		log.Println("Database migration completed successfully!")
	} else {
		log.Fatal("config.DB is nil")
	}
}
