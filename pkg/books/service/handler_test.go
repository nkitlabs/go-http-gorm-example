package service_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/nkitlabs/go-http-gorm-example/pkg/books/service"
	"github.com/nkitlabs/go-http-gorm-example/pkg/books/testutil"
	"github.com/nkitlabs/go-http-gorm-example/pkg/books/types"
	apierrors "github.com/nkitlabs/go-http-gorm-example/pkg/errors"
)

func TestAddBook(t *testing.T) {
	type output struct {
		code   int
		body   types.AddBookResponse
		errMsg string
	}
	testCases := []struct {
		name       string
		preProcess func(s *testutil.TestSuite)
		input      types.AddBookRequest
		output     output
	}{
		{
			name:  "invalid body",
			input: types.AddBookRequest{},
			output: output{
				code:   http.StatusBadRequest,
				body:   types.AddBookResponse{},
				errMsg: "{\"Author\":\"It is required\",\"Description\":\"It is required\",\"Title\":\"It is required\"}",
			},
		},
		{
			name:  "invalid title and description",
			input: types.AddBookRequest{Author: "test"},
			output: output{
				code:   http.StatusBadRequest,
				body:   types.AddBookResponse{},
				errMsg: "{\"Description\":\"It is required\",\"Title\":\"It is required\"}",
			},
		},
		{
			name:  "success",
			input: types.AddBookRequest{Author: "test-author", Title: "test-title", Description: "test-desc"},
			output: output{
				code: http.StatusCreated,
				body: types.AddBookResponse{ID: 1},
			},
			preProcess: func(s *testutil.TestSuite) {
				s.Repository.EXPECT().CreateBook(types.Book{
					Author:      "test-author",
					Title:       "test-title",
					Description: "test-desc",
				}).Return(types.Book{
					ID:          1,
					Author:      "test-author",
					Title:       "test-title",
					Description: "test-desc",
				}, nil)
			},
		},
		{
			name:  "error database",
			input: types.AddBookRequest{Author: "test-author", Title: "test-title", Description: "test-desc"},
			output: output{
				code:   http.StatusInternalServerError,
				body:   types.AddBookResponse{},
				errMsg: "Internal Server Error",
			},
			preProcess: func(s *testutil.TestSuite) {
				s.Repository.EXPECT().CreateBook(types.Book{
					Author:      "test-author",
					Title:       "test-title",
					Description: "test-desc",
				}).Return(types.Book{}, errors.New("database error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := testutil.NewTestSuite(t)

			if tc.preProcess != nil {
				tc.preProcess(&s)
			}

			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(tc.input)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/books", &body)
			require.NoError(t, err)

			resp := httptest.NewRecorder()
			router := http.NewServeMux()
			router = service.InitializeRoutes(router, *s.Handler)
			router.ServeHTTP(resp, req)

			require.Equal(t, tc.output.code, resp.Code)

			if tc.output.code == http.StatusCreated {
				var res types.AddBookResponse
				err := json.NewDecoder(resp.Body).Decode(&res)
				require.NoError(t, err)

				require.Equal(t, tc.output.body, res)
			} else {
				var res apierrors.Error
				err := json.NewDecoder(resp.Body).Decode(&res)
				require.NoError(t, err)

				require.Equal(t, tc.output.errMsg, res.Message)
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	type output struct {
		code   int
		body   types.DeleteBookResponse
		errMsg string
	}
	testCases := []struct {
		name       string
		preProcess func(s *testutil.TestSuite)
		pathID     string
		output     output
	}{
		{
			name:   "invalid id",
			pathID: "not-an-int",
			output: output{
				code:   http.StatusBadRequest,
				body:   types.DeleteBookResponse{},
				errMsg: "invalid id: not-an-int",
			},
		},
		{
			name:   "id not found",
			pathID: "1",
			output: output{
				code:   http.StatusNotFound,
				body:   types.DeleteBookResponse{},
				errMsg: "book not found",
			},
			preProcess: func(s *testutil.TestSuite) {
				s.Repository.EXPECT().GetBook(1).Return(&types.Book{}, gorm.ErrRecordNotFound)
			},
		},
		{
			name:   "success",
			pathID: "1",
			output: output{
				code: http.StatusOK,
				body: types.DeleteBookResponse{},
			},
			preProcess: func(s *testutil.TestSuite) {
				s.Repository.EXPECT().GetBook(1).Return(&types.Book{
					ID:          1,
					Author:      "test-author",
					Title:       "test-title",
					Description: "test-desc",
				}, nil)
				s.Repository.EXPECT().DeleteBook(&types.Book{
					ID:          1,
					Author:      "test-author",
					Title:       "test-title",
					Description: "test-desc",
				}).Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := testutil.NewTestSuite(t)

			if tc.preProcess != nil {
				tc.preProcess(&s)
			}

			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/books/%s", tc.pathID), nil)
			require.NoError(t, err)

			resp := httptest.NewRecorder()
			router := http.NewServeMux()
			router = service.InitializeRoutes(router, *s.Handler)
			router.ServeHTTP(resp, req)

			require.Equal(t, tc.output.code, resp.Code)

			if tc.output.code == http.StatusOK {
				var res types.DeleteBookResponse
				err := json.NewDecoder(resp.Body).Decode(&res)
				require.NoError(t, err)

				require.Equal(t, tc.output.body, res)
			} else {
				var res apierrors.Error
				err := json.NewDecoder(resp.Body).Decode(&res)
				require.NoError(t, err)

				require.Equal(t, tc.output.errMsg, res.Message)
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	type output struct {
		code   int
		body   types.Book
		errMsg string
	}
	testCases := []struct {
		name       string
		preProcess func(s *testutil.TestSuite)
		pathID     string
		input      types.UpdateBookRequest
		output     output
	}{
		{
			name:   "invalid id",
			pathID: "not-an-int",
			output: output{
				code:   http.StatusBadRequest,
				body:   types.Book{},
				errMsg: "invalid id: not-an-int",
			},
		},
		{
			name:   "update only title",
			pathID: "1",
			input: types.UpdateBookRequest{
				Title: "new-title",
			},
			output: output{
				code: http.StatusOK,
				body: types.Book{
					ID: 1, Author: "test-author", Title: "new-title", Description: "test-desc",
				},
			},
			preProcess: func(s *testutil.TestSuite) {
				s.Repository.EXPECT().GetBook(1).Return(&types.Book{
					ID: 1, Author: "test-author", Title: "test-title", Description: "test-desc",
				}, nil)
				s.Repository.EXPECT().UpdateBook(&types.Book{
					ID: 1, Author: "test-author", Title: "new-title", Description: "test-desc",
				}).Return(nil)
				s.Repository.EXPECT().GetBook(1).Return(&types.Book{
					ID: 1, Author: "test-author", Title: "new-title", Description: "test-desc",
				}, nil)
			},
		},
		{
			name:   "database error",
			pathID: "1",
			input: types.UpdateBookRequest{
				Title: "new-title",
			},
			output: output{
				code:   http.StatusInternalServerError,
				errMsg: "Internal Server Error",
			},
			preProcess: func(s *testutil.TestSuite) {
				s.Repository.EXPECT().GetBook(1).Return(&types.Book{
					ID: 1, Author: "test-author", Title: "test-title", Description: "test-desc",
				}, nil)
				s.Repository.EXPECT().UpdateBook(&types.Book{
					ID: 1, Author: "test-author", Title: "new-title", Description: "test-desc",
				}).Return(errors.New("database error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := testutil.NewTestSuite(t)

			if tc.preProcess != nil {
				tc.preProcess(&s)
			}

			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(tc.input)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/books/%s", tc.pathID), &body)
			require.NoError(t, err)

			resp := httptest.NewRecorder()
			router := http.NewServeMux()
			router = service.InitializeRoutes(router, *s.Handler)
			router.ServeHTTP(resp, req)

			require.Equal(t, tc.output.code, resp.Code)

			if tc.output.code == http.StatusOK {
				var res types.Book
				err := json.NewDecoder(resp.Body).Decode(&res)
				require.NoError(t, err)

				require.Equal(t, tc.output.body, res)
			} else {
				var res apierrors.Error
				err := json.NewDecoder(resp.Body).Decode(&res)
				require.NoError(t, err)

				require.Equal(t, tc.output.errMsg, res.Message)
			}
		})
	}
}
