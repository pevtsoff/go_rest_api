package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"rest_api/config"
	"rest_api/controllers"
	"rest_api/models"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
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

	// Generate Swagger JSON at startup from annotations (no local docs folder)
	swaggerJSON, genErr := generateSwaggerJSON()
	if genErr != nil {
		logger.Printf("warning: failed to generate swagger: %v", genErr)
	}

	// Serve the generated Swagger JSON from memory
	engine.GET("/openapi.json", func(c *gin.Context) {
		if swaggerJSON == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "swagger not available"})
			return
		}
		c.Data(http.StatusOK, "application/json", swaggerJSON)
	})

	// Serve Swagger UI pointing to the local generated spec
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/openapi.json")))

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

// generateSwaggerJSON builds the OpenAPI spec from code annotations at runtime.
func generateSwaggerJSON() ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	parser := swag.New(
		swag.SetParseDependency(int(swag.ParseAll)),
		swag.ParseUsingGoList(true),
		swag.SetStrict(false),
	)
	if err := parser.ParseAPIMultiSearchDir([]string{wd}, "main.go", 1); err != nil {
		return nil, err
	}

	return json.Marshal(parser.GetSwagger())
}
