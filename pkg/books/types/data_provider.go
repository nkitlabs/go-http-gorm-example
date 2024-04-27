package types

import "github.com/nkitlabs/go-http-gorm-example/pkg/db"

// DataProvider is the interface for the data provider for a books service
type DataProvider interface {
	CreateBook(book Book) (Book, error)
	UpdateBook(book *Book) error
	DeleteBook(book *Book) error

	GetBooks(page int, limit int, sortType db.SortType) (*db.Pagination, []Book, error)
	GetBook(id int) (*Book, error)
}
