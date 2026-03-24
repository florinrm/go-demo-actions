package usecases

import (
	"errors"
	"fmt"
	"rest-api/domain"

	"github.com/google/uuid"
)

// DummyDatabase represents a dummy in-memory database
type DummyDatabase struct {
	storage map[string]domain.Book
}

// NewDatabase creates a new instance of DummyDatabase
func NewDatabase() *DummyDatabase {
	return &DummyDatabase{make(map[string]domain.Book)}
}

// Add generates a UUID for the book to be added in database
func (db *DummyDatabase) Add(book domain.Book) string {
	id := uuid.New().String()
	db.storage[id] = book

	return id
}

// Update is used for updating an existing book
func (db *DummyDatabase) Update(hash string, book domain.Book) {
	db.storage[hash] = book
}

// RemoveByID removes a book from database using its ID, if found
func (db *DummyDatabase) RemoveByID(id string) error {
	_, ok := db.storage[id]
	if !ok {
		return errors.New(fmt.Sprintf("book with id %s not found", id))
	}

	delete(db.storage, id)
	return nil
}

// GetByID retrieves a book by its ID, if found
func (db *DummyDatabase) GetByID(id string) (domain.Book, error) {
	book, ok := db.storage[id]
	if !ok {
		return domain.Book{}, errors.New("book not found")
	}

	return book, nil
}

// GetAll returns all books saved in database
func (db *DummyDatabase) GetAll(filters map[string]interface{}) []domain.Book {
	var books []domain.Book
	for _, book := range db.storage {
		if db.checkFilters(book, filters) {
			books = append(books, book)
		}
	}
	return books
}

func (db *DummyDatabase) checkFilters(book domain.Book, filters map[string]interface{}) bool {
	if author, ok := filters["author"]; ok {
		if book.Author != author {
			return false
		}
	}

	if title, ok := filters["title"]; ok {
		if book.Title != title {
			return false
		}
	}

	if genre, ok := filters["genre"]; ok {
		if book.Genre != genre {
			return false
		}
	}

	if year, ok := filters["year"]; ok {
		if book.Year != year {
			return false
		}
	}

	return true
}
