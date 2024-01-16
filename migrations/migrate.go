package main

import (
	"fmt"
	"log"
	"os"
	"rest_api/config"
	"rest_api/models"
)


func printEnvVars(){
	envVars := os.Environ()

	for _, envVar := range envVars {
		fmt.Println(envVar)
	}
}

func init(){
	config.LoadEnvVars()
	printEnvVars()
	config.ConnectToDB()
}


func main(){
	if config.DB != nil {
		config.DB.AutoMigrate(&models.Post{})
	} else {
		log.Fatal("config.DB is nil")
	}
}
