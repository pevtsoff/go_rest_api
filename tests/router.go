package tests

import (
	"rest_api/controllers"
	errors_middleware "rest_api/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter builds a gin.Engine with the same routes and middleware as main.go
func NewRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(errors_middleware.JSONErrorMiddleware())

	r.POST("/users/", controllers.UsersCreate)
	r.GET("/users/:id", controllers.UsersShow)
	r.GET("/users/:id/posts", controllers.UserPostsShow)
	r.POST("/posts/", controllers.PostsCreate)
	r.GET("/posts/", controllers.PostsIndex)
	r.GET("/posts/:id", controllers.PostsShow)
	r.PATCH("/posts/:id", controllers.PostsUpdate)
	r.DELETE("/posts/:id", controllers.PostsDelete)

	return r
}
