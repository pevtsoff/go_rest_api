package controllers

import (
	"errors"
	"rest_api/config"
	"rest_api/models"
	"time"

	"github.com/gin-gonic/gin"
)


// CreatePostRequest represents the request body for creating a post
type CreatePostRequest struct {
	Title string `json:"title" binding:"required" example:"My First Post"`
	Body  string `json:"body" binding:"required" example:"This is the content of my first post"`
}

// UpdatePostRequest represents the request body for updating a post
type UpdatePostRequest struct {
	Title string `json:"title" example:"Updated Post Title"`
	Body  string `json:"body" example:"Updated post content"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name string `json:"name" binding:"required" example:"John Doe"`
}

// PostsResponse represents the response for posts endpoints
type PostsResponse struct {
	Posts []models.JsonPost `json:"posts"`
}

// PostResponse represents the response for a single post
type PostResponse struct {
	Post models.JsonPost `json:"post"`
}

// UsersResponse represents the response for users endpoints
type UsersResponse struct {
	Users []models.JsonUser `json:"users"`
}

// UserResponse represents the response for a single user
type UserResponse struct {
	User models.JsonUser `json:"user"`
}

// mapPost converts DB model to API DTO
func mapPost(m models.Post) models.JsonPost {
	var deletedAt *string
	if m.DeletedAt.Valid {
		s := m.DeletedAt.Time.Format(time.RFC3339)
		deletedAt = &s
	}
	return models.JsonPost{
		ID:        m.ID,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
		DeletedAt: deletedAt,
		Title:     m.Title,
		Body:      m.Body,
	}
}

// mapUser converts DB model to API DTO
func mapUser(m models.User) models.JsonUser {
	var deletedAt *string
	if m.DeletedAt.Valid {
		s := m.DeletedAt.Time.Format(time.RFC3339)
		deletedAt = &s
	}
	return models.JsonUser{
		ID:        m.ID,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
		DeletedAt: deletedAt,
		Name:      m.Name,
	}
}

// PostsCreate godoc
// @Summary Create a new post
// @Description Create a new blog post
// @Tags posts
// @Accept json
// @Produce json
// @Param post body CreatePostRequest true "Post data"
// @Success 200 {object} PostResponse "Post created successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Router /posts [post]
func PostsCreate(c *gin.Context) {
	var body CreatePostRequest
	err := c.Bind(&body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	post := models.Post{Title: body.Title, Body: body.Body}
	result := config.DB.Create(&post)

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(200, PostResponse{Post: mapPost(post)})
}

// PostsIndex godoc
// @Summary Get all posts
// @Description Get a list of all blog posts
// @Tags posts
// @Produce json
// @Success 200 {object} PostsResponse "List of posts"
// @Router /posts [get]
func PostsIndex(c *gin.Context) {
	var posts []models.Post
	config.DB.Find(&posts)

	dto := make([]models.JsonPost, 0, len(posts))
	for _, p := range posts {
		dto = append(dto, mapPost(p))
	}
	c.JSON(200, PostsResponse{Posts: dto})
}

// PostsShow godoc
// @Summary Get a post by ID
// @Description Get a specific blog post by its ID
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} PostResponse "Post found"
// @Failure 404 {object} map[string]string "Post not found"
// @Router /posts/{id} [get]
func PostsShow(c *gin.Context) {
	var post models.Post
	id := c.Param("id")
	result := config.DB.First(&post, id)

	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(200, PostResponse{Post: mapPost(post)})
}

// PostsUpdate godoc
// @Summary Update a post
// @Description Update an existing blog post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body UpdatePostRequest true "Updated post data"
// @Success 200 {object} PostResponse "Post updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Post not found"
// @Router /posts/{id} [patch]
func PostsUpdate(c *gin.Context) {
	var body UpdatePostRequest
	err := c.Bind(&body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var post models.Post
	result := config.DB.First(&post, c.Param("id"))

	if result.Error != nil {
		c.Error(errors.New("Unable to update a post"))
		return
	}

	config.DB.Model(&post).Updates(models.Post{
		Title: body.Title,
		Body:  body.Body,
	})

	c.JSON(200, PostResponse{Post: mapPost(post)})
}

// PostsDelete godoc
// @Summary Delete a post
// @Description Delete a blog post by ID
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]string "Post deleted successfully"
// @Failure 404 {object} map[string]string "Post not found"
// @Router /posts/{id} [delete]
func PostsDelete(c *gin.Context) {
	var post models.Post
	result := config.DB.First(&post, c.Param("id"))

	if result.Error != nil {
		c.Error(errors.New("Unable to delete a post"))
		return
	}

	config.DB.Model(&post).Delete(&post)

	c.JSON(
		200, gin.H{
			"id has been deleted": c.Param("id"),
		})
}

// UsersCreate godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User data"
// @Success 200 {object} UserResponse "User created successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Router /users [post]
func UsersCreate(c *gin.Context) {
	var body CreateUserRequest
	err := c.Bind(&body)
	if err != nil {
		c.Error(errors.New(err.Error()))
		return
	}

	user := models.User{Name: body.Name}
	result := config.DB.Create(&user)

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(200, UserResponse{User: mapUser(user)})
}

// UsersShow godoc
// @Summary Get a user by ID
// @Description Get a specific user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse "User found"
// @Failure 404 {object} map[string]string "User not found"
// @Router /users/{id} [get]
func UsersShow(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	result := config.DB.First(&user, id)

	if result.Error != nil {
		c.Error(errors.New("User not found"))
		return
	}

	c.JSON(200, UserResponse{User: mapUser(user)})
}
