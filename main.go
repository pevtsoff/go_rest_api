package main

import (
	"log"
	"rest_api/config"
	"rest_api/controllers"
	"rest_api/models"

	_ "rest_api/docs" // This is required for swagger to work

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

// @title REST API
// @version 1.0
// @description A simple REST API for managing posts and users
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /
func main() {
	logger.Println("hello world")
	engine := gin.Default()
	engine.Use(JSONErrorMiddleware())

	// Swagger documentation endpoint
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	engine.POST("/users/", controllers.UsersCreate)
	engine.GET("/users/:id", controllers.UsersShow)
	engine.POST("/posts/", controllers.PostsCreate)
	engine.GET("/posts/", controllers.PostsIndex)
	engine.GET("/posts/:id", controllers.PostsShow)
	engine.PATCH("/posts/:id", controllers.PostsUpdate)
	engine.DELETE("/posts/:id", controllers.PostsDelete)

	engine.Run(":3000") // listen and serve on localhost:3000
}
