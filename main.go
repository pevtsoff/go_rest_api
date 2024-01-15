package main

import (
  "log"
  "net/http"
  "github.com/gin-gonic/gin"
   "github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error: Cannot load .env file")
	}
}

func main() {
  print("hello world")
  engine := gin.Default()
  engine.GET("/", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "hello world",
    })
  })
  engine.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}