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

func PostsIndex(c *gin.Context){
	var posts []models.Post
	config.DB.Find(&posts)

	c.JSON(
	200, gin.H{
		"post": posts,
	})
}


func PostsShow(c *gin.Context){
	var post models.Post
	id := c.Param("id")
	config.DB.Find(&post, id)

	c.JSON(
	200, gin.H{
		"post": post,
	})
}


func PostsUpdate(c *gin.Context){
	var body struct{
		Title string
		Body string
	}
	err := c.Bind(&body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var post  models.Post
    result := config.DB.First(&post, c.Param("id"))

	if result.Error !=nil {
		c.Status(400)
		return
	}

	config.DB.Model(&post).Updates(models.Post{
		Title: body.Title,
		Body: body.Body,
	})
	
	c.JSON(
	200, gin.H{
		"post": post,
	})
}
