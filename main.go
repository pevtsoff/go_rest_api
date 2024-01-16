package main

import (
	"log"
	"net/http"
	"rest_api/config"
	"github.com/gin-gonic/gin"
)

var logger = log.Default()

func init() {
	config.LoadEnvVars()
  config.ConnectToDB()
}

func main() {
  logger.Println("hello world")
  engine := gin.Default()
  engine.GET("/", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "hello world",
    })
  })
  engine.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}