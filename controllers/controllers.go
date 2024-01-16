package controllers

import (
	"rest_api/config"
	"rest_api/models"

	"github.com/gin-gonic/gin"
)

func PostsCreate(c *gin.Context){
	var body struct{
		Title string
		Body string
	}
	err := c.Bind(&body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	post := models.Post{Title: body.Title, Body: body.Body}
    result := config.DB.Create(&post)

	if result.Error !=nil {
		c.Status(400)
		return
	}

	c.JSON(
	200, gin.H{
		"post": post,
	})
}