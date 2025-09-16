package models
// Post represents a blog post
type JsonPost struct {
	ID        uint    `json:"id" example:"1"`
	CreatedAt string  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt string  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt *string `json:"deleted_at,omitempty" example:"null"`
	Title     string  `json:"title" example:"My First Post"`
	Body      string  `json:"body" example:"This is the content of my first post"`
}

// User represents a user
type JsonUser struct {
	ID        uint    `json:"id" example:"1"`
	CreatedAt string  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt string  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt *string `json:"deleted_at,omitempty" example:"null"`
	Name      string  `json:"name" example:"John Doe"`
}
