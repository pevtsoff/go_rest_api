package main

import (
	"log"
	"rest_api/config"
	"rest_api/controllers"

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
  engine.POST("/posts/", controllers.PostsCreate)
  engine.GET("/posts/", controllers.PostsIndex)
  engine.GET("/posts/:id", controllers.PostsShow)
  engine.PATCH("/posts/:id", controllers.PostsUpdate)
  engine.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}