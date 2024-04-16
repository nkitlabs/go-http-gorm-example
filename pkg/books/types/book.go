package types

// Book is the model for a book information.
type Book struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`
	Title       string `json:"title" example:"example-title" validate:"required"`
	Author      string `json:"author" example:"John Doe" validate:"required"`
	Description string `json:"description" example:"this is an example description" validate:"required"`
}
