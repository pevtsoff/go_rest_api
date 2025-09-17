package testutils

import (
	"time"

	"rest_api/config"
	"rest_api/models"
)

// UserBuilder helps create a User row for tests and returns the created model.
type UserBuilder struct {
	name string
}

// New initializes (or re-initializes) the builder with default values.
func (b *UserBuilder) New() *UserBuilder {
	b.name = "Test User " + time.Now().UTC().Format(time.RFC3339Nano)
	return b
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
	b.name = name
	return b
}

// Create inserts the user into DB using the current config.DB (can be a tx) and returns it.
func (b *UserBuilder) Create() (models.User, error) {
	u := models.User{Name: b.name}
	if err := config.DB.Create(&u).Error; err != nil {
		return models.User{}, err
	}
	return u, nil
}

// PostBuilder helps create a Post row for tests and returns the created model.
type PostBuilder struct {
	title string
	body  string
}

// New initializes (or re-initializes) the builder with default values.
func (b *PostBuilder) New() *PostBuilder {
	ts := time.Now().UTC().Format(time.RFC3339Nano)
	b.title = "Test Title " + ts
	b.body = "Test Body " + ts
	return b
}

func (b *PostBuilder) WithTitle(title string) *PostBuilder {
	b.title = title
	return b
}

func (b *PostBuilder) WithBody(body string) *PostBuilder {
	b.body = body
	return b
}

// Create inserts the post into DB using the current config.DB (can be a tx) and returns it.
func (b *PostBuilder) Create() (models.Post, error) {
	p := models.Post{Title: b.title, Body: b.body}
	if err := config.DB.Create(&p).Error; err != nil {
		return models.Post{}, err
	}
	return p, nil
}

// Itoa converts a uint ID to string for test URL building without importing strconv.
func Itoa(u uint) string {
	if u == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for u > 0 {
		i--
		buf[i] = byte('0' + u%10)
		u /= 10
	}
	return string(buf[i:])
}
