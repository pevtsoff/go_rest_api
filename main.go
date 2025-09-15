package main

import (
	"log"
	"rest_api/config"
	"rest_api/controllers"
	"rest_api/models"

	"github.com/gin-gonic/gin"
)

var logger = log.Default()

func init() {
	config.LoadEnvVars()
	config.ConnectToDB()

	// Auto-migrate the database schema
	err := config.DB.AutoMigrate(&models.Post{}, &models.User{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Database migration completed successfully!")
}

func JSONErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.AbortWithStatusJSON(c.Writer.Status(), gin.H{
				"error": c.Errors[0].Error(),
			})
		}
	}
}

func main() {
	logger.Println("hello world")
	engine := gin.Default()
	engine.Use(JSONErrorMiddleware())
	engine.POST("/users/", controllers.UsersCreate)
	engine.GET("/users/:id", controllers.UsersShow)
	engine.POST("/posts/", controllers.PostsCreate)
	engine.GET("/posts/", controllers.PostsIndex)
	engine.GET("/posts/:id", controllers.PostsShow)
	engine.PATCH("/posts/:id", controllers.PostsUpdate)
	engine.DELETE("/posts/:id", controllers.PostsDelete)
	engine.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
