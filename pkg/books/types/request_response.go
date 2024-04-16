package types

import "github.com/nkitlabs/go-http-gorm-example/pkg/db"

// AddBookRequest is the request for adding a new book
type AddBookRequest struct {
	Title       string `json:"title" example:"example-title" validate:"required"`
	Author      string `json:"author" example:"John Doe" validate:"required"`
	Description string `json:"description" example:"this is an example description" validate:"required"`
}

// AddBookResponse is the response for adding a new book
type AddBookResponse struct {
	ID int `json:"id" example:"1"`
}

// GetBooksResponse is the response for getting a list of books
type GetBooksResponse struct {
	Books      []Book         `json:"books"`
	Pagination *db.Pagination `json:"pagination"`
}

// DeleteBookResponse is the response for deleting a book
type DeleteBookResponse struct{}

// UpdateBookRequest is the request for updating a book information
type UpdateBookRequest struct {
	Title       string `json:"title" example:"example-title"`
	Author      string `json:"author" example:"John Doe"`
	Description string `json:"description" example:"this is an example description"`
}
