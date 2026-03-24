package gateways

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"rest-api/domain"
	"rest-api/mocks"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mocksAPITest struct {
	mockDatabase mocks.Database
}

func initMocksAPI() (*BooksAPI, *mocksAPITest) {
	var mocksAPI mocksAPITest
	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	api := NewBooksAPI(&mocksAPI.mockDatabase, &logger)
	api.Handle()
	return api, &mocksAPI
}

func (m *mocksAPITest) assertExpectations(t *testing.T) {
	m.mockDatabase.AssertExpectations(t)
}

func TestBooksAPI_AddBook_Success(t *testing.T) {
	api, apiMocks := initMocksAPI()

	apiMocks.mockDatabase.On("Add", mock.Anything).Return("1234")

	book := domain.Book{
		Title:  "1984",
		Author: "George Orwell",
		Genre:  "Novel",
		Year:   1949,
	}

	reqBytes, _ := json.Marshal(book)

	req, _ := http.NewRequest(http.MethodPost, "/entities/books", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.GetRouter().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response domain.Response
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, response.Book, book)

	apiMocks.assertExpectations(t)
}

func TestBooksAPI_GetBook_Success(t *testing.T) {
	api, apiMocks := initMocksAPI()

	book := domain.Book{
		Title:  "1984",
		Author: "George Orwell",
		Genre:  "Novel",
		Year:   1949,
	}
	id := "1234"
	apiMocks.mockDatabase.On("GetByID", mock.Anything).Return(book, nil)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/entities/books/%s", id), nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.GetRouter().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response domain.Response
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, response.Book, book)

	apiMocks.assertExpectations(t)
}

func TestBooksAPI_GetBook_NotFound(t *testing.T) {
	api, apiMocks := initMocksAPI()

	id := "1234"
	apiMocks.mockDatabase.On("GetByID", mock.Anything).Return(domain.Book{}, errors.New("not found"))

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/entities/books/%s", id), nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.GetRouter().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	apiMocks.assertExpectations(t)
}

func TestBooksAPI_GetAll_Success(t *testing.T) {
	api, apiMocks := initMocksAPI()

	book1 := domain.Book{
		Title:  "1984",
		Author: "George Orwell",
		Genre:  "Novel",
		Year:   1949,
	}

	book2 := domain.Book{
		Title:  "Animal Farm",
		Author: "George Orwell",
		Genre:  "Novella",
		Year:   1945,
	}

	apiMocks.mockDatabase.On("GetAll", mock.Anything).Return([]domain.Book{book1, book2}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/queries/books", nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.GetRouter().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var books []domain.Book
	err := json.Unmarshal(rr.Body.Bytes(), &books)
	assert.Nil(t, err)
	assert.Equal(t, len(books), 2)

	apiMocks.assertExpectations(t)
}

func TestBooksAPI_Delete_Success(t *testing.T) {
	api, apiMocks := initMocksAPI()

	id := "1234"
	apiMocks.mockDatabase.On("RemoveByID", id).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/entities/books/%s", id), nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.GetRouter().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	apiMocks.assertExpectations(t)
}

func TestBooksAPI_Delete_NotFound(t *testing.T) {
	api, apiMocks := initMocksAPI()

	id := "1234"
	apiMocks.mockDatabase.On("RemoveByID", id).Return(errors.New("not found"))

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/entities/books/%s", id), nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.GetRouter().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	apiMocks.assertExpectations(t)
}

func TestBooksAPI_Put_Success(t *testing.T) {
	api, apiMocks := initMocksAPI()

	apiMocks.mockDatabase.On("Update", mock.Anything, mock.Anything)

	book := domain.Book{
		Title:  "1984",
		Author: "George Orwell",
		Genre:  "Novel",
		Year:   1949,
	}
	id := "1234"

	reqBytes, _ := json.Marshal(book)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/entities/books/%s", id), bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.GetRouter().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response domain.Response
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, response.Book, book)

	apiMocks.assertExpectations(t)
}

func TestBooksAPI_Patch_Success(t *testing.T) {
	api, apiMocks := initMocksAPI()

	apiMocks.mockDatabase.On("GetByID", mock.Anything).Return(domain.Book{
		Title: "1984",
	}, nil)
	apiMocks.mockDatabase.On("Update", mock.Anything, mock.Anything)

	book := domain.Book{
		Title:  "1984",
		Author: "George Orwell",
		Genre:  "Novel",
		Year:   1949,
	}
	id := "1234"

	reqBytes, _ := json.Marshal(book)

	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/entities/books/%s", id), bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.GetRouter().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response domain.Response
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, response.Book, book)

	apiMocks.assertExpectations(t)
}

func TestBooksAPI_Patch_NotFound(t *testing.T) {
	api, apiMocks := initMocksAPI()

	apiMocks.mockDatabase.On("GetByID", mock.Anything).Return(domain.Book{}, errors.New("not found"))

	book := domain.Book{
		Title:  "1984",
		Author: "George Orwell",
		Genre:  "Novel",
		Year:   1949,
	}
	id := "1234"

	reqBytes, _ := json.Marshal(book)

	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/entities/books/%s", id), bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.GetRouter().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	apiMocks.assertExpectations(t)
}
