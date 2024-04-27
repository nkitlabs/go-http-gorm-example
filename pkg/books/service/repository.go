package service

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/nkitlabs/go-http-gorm-example/pkg/books/types"
	"github.com/nkitlabs/go-http-gorm-example/pkg/db"
)

var (
	_ types.DataProvider = &Repository{}
)

// Repository is the data provider that connect to the database being used in a books service.
type Repository struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewRepository creates a new books-service repository
func NewRepository(db *gorm.DB, log *zap.Logger) Repository {
	return Repository{db, log}
}

// CreateBook creates a new book in the database
func (r *Repository) CreateBook(book types.Book) (types.Book, error) {
	tx := r.db.Create(&book)
	return book, tx.Error
}

// UpdateBook updates a book in the database
func (r *Repository) UpdateBook(book *types.Book) error {
	return r.db.Save(&book).Error
}

// DeleteBook deletes a book from the database
func (r *Repository) DeleteBook(book *types.Book) error {
	return r.db.Delete(book).Error
}

// GetBooks retrieves a list of books from the database
func (r *Repository) GetBooks(page int, limit int, sortType db.SortType) (*db.Pagination, []types.Book, error) {
	p := db.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  fmt.Sprintf("id %s", sortType),
	}

	var books []types.Book
	if result := r.db.Scopes(db.Paginate(&books, &p, r.db)).Find(&books); result.Error != nil {
		return nil, nil, result.Error
	}

	return &p, books, nil
}

// GetBook retrieves a book from the database
func (r *Repository) GetBook(id int) (*types.Book, error) {
	var book types.Book
	if result := r.db.First(&book, id); result.Error != nil {
		return nil, result.Error
	}

	return &book, nil
}
