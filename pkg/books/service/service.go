package service

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/nkitlabs/go-http-gorm-example/pkg/books/types"
	"github.com/nkitlabs/go-http-gorm-example/pkg/db"
	apierror "github.com/nkitlabs/go-http-gorm-example/pkg/errors"
)

// Service is the service layer for books
type Service struct {
	dataProvider types.DataProvider
	log          *zap.Logger
}

// NewService creates a new books service
func NewService(d types.DataProvider, log *zap.Logger) Service {
	return Service{
		dataProvider: d,
		log:          log,
	}
}

// AddBook adds a new book into a system
func (s *Service) AddBook(req types.AddBookRequest) (types.AddBookResponse, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return types.AddBookResponse{}, apierror.ConvertValidatorErrorsToError(err)
	}

	book := types.Book{
		Title:       req.Title,
		Author:      req.Author,
		Description: req.Description,
	}

	book, err := s.dataProvider.CreateBook(book)
	if err != nil {
		return types.AddBookResponse{}, err
	}

	return types.AddBookResponse{
		ID: book.ID,
	}, nil
}

// UpdateBook updates a book information in the system
func (s *Service) UpdateBook(id int, req types.UpdateBookRequest) (types.Book, error) {
	book, err := s.GetBook(id)
	if err != nil {
		return types.Book{}, err
	}

	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.Description != "" {
		book.Description = req.Description
	}

	if err := s.dataProvider.UpdateBook(book); err != nil {
		return types.Book{}, err
	}

	newBook, err := s.GetBook(id)
	if err != nil {
		return types.Book{}, err
	}

	return *newBook, nil
}

// DeleteBook deletes a book id from the system
func (s *Service) DeleteBook(id int) (types.DeleteBookResponse, error) {
	book, err := s.GetBook(id)
	if err != nil {
		return types.DeleteBookResponse{}, err
	}

	return types.DeleteBookResponse{}, s.dataProvider.DeleteBook(book)
}

// GetBooks returns a list of books
func (s *Service) GetBooks(page int, limit int, sortType db.SortType) (types.GetBooksResponse, error) {
	pagination, books, err := s.dataProvider.GetBooks(page, limit, sortType)
	if err != nil {
		return types.GetBooksResponse{}, err
	}

	return types.GetBooksResponse{
		Books:      books,
		Pagination: pagination,
	}, nil
}

// GetBook returns a book information from the given id
func (s *Service) GetBook(id int) (*types.Book, error) {
	book, err := s.dataProvider.GetBook(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apierror.NewNotFoundError("book not found")
	} else if err != nil {
		return nil, err
	}

	return book, nil
}
