package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	_ "github.com/nkitlabs/go-http-gorm-example/docs"
	"github.com/nkitlabs/go-http-gorm-example/pkg/books/types"
	"github.com/nkitlabs/go-http-gorm-example/pkg/db"
	apierror "github.com/nkitlabs/go-http-gorm-example/pkg/errors"
	"github.com/nkitlabs/go-http-gorm-example/pkg/response"
)

type Handler struct {
	serv *Service
	log  *zap.Logger
}

// New creates a new handler for the books service
func NewHandler(serv *Service, log *zap.Logger) Handler {
	return Handler{serv, log}
}

// InitializeRoutes initializes the routes for the books service
func InitializeRoutes(mux *http.ServeMux, h Handler) *http.ServeMux {
	mux.HandleFunc("GET /api/v1/books", h.GetBooks)
	mux.HandleFunc("GET /api/v1/books/{id}", h.GetBook)
	mux.HandleFunc("POST /api/v1/books", h.AddBook)
	mux.HandleFunc("PUT /api/v1/books/{id}", h.UpdateBook)
	mux.HandleFunc("DELETE /api/v1/books/{id}", h.DeleteBook)
	return mux
}

// @Summary get book information from given id
// @ID get-book
// @Param id path int true "Book ID"
// @Produce json
// @Success 200 {object} types.Book
// @Failure 404 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Router /books/{id} [get]
func (h Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	// Read the dynamic id parameter
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 0)
	if err != nil {
		response.WriteError(w, err, h.log)
		return
	}

	// Find book by Id
	book, err := h.serv.GetBook(int(id))
	if err != nil {
		response.WriteError(w, err, h.log)
		return
	}

	response.Write(w, http.StatusOK, book, h.log)
}

// @Summary get list of books' information from the system
// @ID get-books
// @Param page query int true "Page number"
// @Param limit query int true "Limit per page"
// @Param sort_type query string false "order of items to be sorted (by id)" enums(asc,desc)
// @Produce json
// @Success 200 {object} types.GetBooksResponse
// @Failure 500 {object} errors.Error
// @Router /books [get]
func (h Handler) GetBooks(w http.ResponseWriter, r *http.Request) {
	// Get page and limit from params and convert to int
	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 0)
	if err != nil {
		response.WriteError(w, err, h.log)
		return
	}

	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 0)
	if err != nil {
		response.WriteError(w, err, h.log)
		return
	}

	sortType := db.ToSortType(r.URL.Query().Get("sort_type"))

	result, err := h.serv.GetBooks(int(page), int(limit), sortType)
	if err != nil {
		response.WriteError(w, err, h.log)
		return
	}

	response.Write(w, http.StatusOK, result, h.log)
}

// @Summary add book into the system
// @ID add-book
// @Produce json
// @Param Body body types.AddBookRequest true "Book information that needs to be added"
// @Success 201 {object} types.AddBookResponse
// @Failure 500 {object} errors.Error
// @Router /books [post]
func (h Handler) AddBook(w http.ResponseWriter, r *http.Request) {
	// Read to request body
	var req types.AddBookRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, err, h.log)
		return
	}

	// Append to the Books table
	resp, err := h.serv.AddBook(req)
	if err != nil {
		response.WriteError(w, err, h.log)
		return
	}

	response.Write(w, http.StatusCreated, resp, h.log)
}

// @Summary delete book id from the system
// @ID delete-book
// @Param id path int true "Book ID"
// @Produce json
// @Success 200 {object} types.DeleteBookResponse
// @Failure 404 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Router /books/{id} [delete]
func (h Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	// Read the dynamic id parameter
	pathID := r.PathValue("id")
	id, err := strconv.ParseInt(pathID, 10, 0)
	if err != nil {
		h.log.Error(fmt.Sprintf("%s; invalid id %s", err.Error(), pathID))
		newErr := apierror.ErrInvalidInput.WithMessage(fmt.Sprintf("invalid id: %s", r.PathValue("id")))
		response.WriteError(w, newErr, h.log)
		return
	}

	// Find the book by Id
	result, err := h.serv.DeleteBook(int(id))
	if err != nil {
		response.WriteError(w, err, h.log)
		return
	}

	response.Write(w, http.StatusOK, result, h.log)
}

// @Summary update book information in the system with the given id
// @ID add-book
// @Produce json
// @Param id path int true "Book ID"
// @Param Body body types.UpdateBookRequest true "Book information that needs to be updated"
// @Success 200 {object} types.Book
// @Failure 404 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Router /books/{id} [put]
func (h Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	// Read the dynamic id parameter
	pathID := r.PathValue("id")
	id, err := strconv.ParseInt(pathID, 10, 0)
	if err != nil {
		h.log.Error(fmt.Sprintf("%s; invalid id %s", err.Error(), pathID))
		newErr := apierror.ErrInvalidInput.WithMessage(fmt.Sprintf("invalid id: %s", r.PathValue("id")))
		response.WriteError(w, newErr, h.log)
		return
	}

	// Read request body
	defer r.Body.Close()
	var req types.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, err, h.log)
		return
	}

	result, err := h.serv.UpdateBook(int(id), req)
	if err != nil {
		response.WriteError(w, err, h.log)
		return
	}

	response.Write(w, http.StatusOK, result, h.log)
}
